from pathlib import Path
import re

import torch

from utils import get_device, generate_text_basic_stream

from reasoning_from_scratch.qwen3 import (
    download_qwen3_small,
    Qwen3Tokenizer,
    Qwen3Model,
    QWEN_CONFIG_06_B,
    KVCache,
)


def load_model_and_tokenizer(which_model, device, use_compile, local_dir="qwen3"):
    if which_model == "base":
        download_qwen3_small(kind="base", tokenizer_only=False, out_dir=local_dir)
        tokenizer_path = Path(local_dir) / "tokenizer-base.json"
        model_path = Path(local_dir) / "qwen3-0.6B-base.pth"
        tokenizer = Qwen3Tokenizer(tokenizer_file_path=tokenizer_path)
    elif which_model == "reasoning":
        download_qwen3_small(kind="reasoning", tokenizer_only=False, out_dir=local_dir)
        tokenizer_path = Path(local_dir) / "tokenizer-reasoning.json"
        model_path = Path(local_dir) / "qwen3-0.6B-reasoning.pth"
        tokenizer = Qwen3Tokenizer(
            tokenizer_file_path=tokenizer_path,
            apply_chat_template=True,
            add_generation_prompt=True,
            add_thinking=True,
        )
    else:
        raise ValueError(f"Invalid choice: which_model={which_model}")

    model = Qwen3Model(QWEN_CONFIG_06_B)
    model.load_state_dict(torch.load(model_path))

    model.to(device)

    if use_compile:
        torch._dynamo.config.allow_unspec_int_on_nn_module = True
        model = torch.compile(model)

    return model, tokenizer


WHICH_MODEL = "base"
#WHICH_MODEL = "reasoning"
device = get_device()

model, tokenizer = load_model_and_tokenizer(
    which_model=WHICH_MODEL,
    device=device,
    use_compile=False,
)

prompt = (  # A
    r"If $a+b=3$ and $ab=\tfrac{13}{6}$, "
    r"what is the value of $a^2+b^2$?"
)

input_token_ids_tensor = torch.tensor(
    tokenizer.encode(prompt),
    device=device,
).unsqueeze(0)

all_token_ids = []


# for token in generate_text_basic_stream(
#     model=model,
#     token_ids=input_token_ids_tensor,
#     max_new_tokens=2048,
#     eos_token_id=tokenizer.eos_token_id,
# ):
#     token_id = token.squeeze(0)
#     decoded_id = tokenizer.decode(token_id.tolist())
#     print(
#         decoded_id,
#         end="",
#         flush=True,
#     )
#     all_token_ids.append(token_id)

# all_tokens = tokenizer.decode(all_token_ids)


def generate_text_stream_concat(
    model,
    tokenizer,
    prompt,
    device,
    max_new_tokens,
    verbose=False,
):
    input_ids = torch.tensor(tokenizer.encode(prompt), device=device).unsqueeze(0)

    generated_ids = []
    for token in generate_text_basic_stream(
        model=model,
        token_ids=input_ids,
        max_new_tokens=max_new_tokens,
        eos_token_id=tokenizer.eos_token_id,
    ):
        next_token_id = token.squeeze(0)
        generated_ids.append(next_token_id.item())

        if verbose:
            print(
                tokenizer.decode(next_token_id.tolist()),
                end="",
                flush=True,
            )

    return tokenizer.decode(generated_ids)


# generated_text = generate_text_stream_concat(
#     model, tokenizer, prompt, device,
#     max_new_tokens=2048,
#     verbose=True,
# )


model_answer = r"""... some explanation...
**Final Answer:**
\[
\boxed{\dfrac{14}{3}}
\]
"""


def get_last_boxed(text):
    boxed_start_idx = text.rfind(r"\boxed")
    if boxed_start_idx == -1:
        return None

    # position after "\boxed"
    current_idx = boxed_start_idx + len(r"\boxed")

    # skip whitespace
    while current_idx < len(text) and text[current_idx].isspace():
        current_idx += 1

    # next char must be an opening brace
    if current_idx >= len(text) or text[current_idx] != "{":
        return None

    current_idx += 1
    brace_depth = 1
    content_start_idx = current_idx

    # find matching closing brace
    while current_idx < len(text) and brace_depth > 0:
        char = text[current_idx]
        if char == "{":
            brace_depth += 1
        elif char == "}":
            brace_depth -= 1
        current_idx += 1

    if brace_depth != 0:
        return None

    return text[content_start_idx : current_idx - 1]


extracted_answer = get_last_boxed(model_answer)
print(extracted_answer)


RE_NUMBER = re.compile(  # A
    r"-?(?:\d+/\d+|\d+(?:\.\d+)?(?:[eE][+-]?\d+)?)"
)


def extract_final_candidate(text, fallback="number_then_full"):
    result = ""

    if text:
        boxed = get_last_boxed(text.strip())
        if boxed:
            result = boxed.strip().strip("$ ")
        elif fallback in ("number_then_full", "number_only"):
            m = RE_NUMBER.findall(text)
            if m:
                result = m[-1]
            elif fallback == "number_then_full":
                result = text
    return result


print(extract_final_candidate(model_answer))


LATEX_FIXES = [  # A
    (r"\\left\s*", ""),
    (r"\\right\s*", ""),
    (r"\\,|\\!|\\;|\\:", ""),
    (r"\\cdot", "*"),
    (r"\u00B7|\u00D7", "*"),
    (r"\\\^\\circ", ""),
    (r"\\dfrac", r"\\frac"),
    (r"\\tfrac", r"\\frac"),
    (r"°", ""),
]

RE_SPECIAL = re.compile(r"<\|[^>]+?\|>")  # B
SUPERSCRIPT_MAP = {
    "⁰": "0",
    "¹": "1",
    "²": "2",
    "³": "3",
    "⁴": "4",
    "⁵": "5",
    "⁶": "6",
    "⁷": "7",
    "⁸": "8",
    "⁹": "9",
    "⁺": "+",
    "⁻": "-",
    "⁽": "(",
    "⁾": ")",
}


def normalize_text(text):
    if not text:
        return ""

    # remove special token-like substrings
    text = RE_SPECIAL.sub("", text).strip()

    # remove leading single-letter labels like "A. ..." or "a: ..."
    match = re.match(r"^[A-Za-z]\s*[.:]\s*(.+)$", text)
    if match:
        text = match.group(1)

    # remove degree superscript patterns
    text = re.sub(r"\^\s*\{\s*\\circ\s*\}", "", text)
    text = re.sub(r"\^\s*\\circ", "", text)

    # remove explicit degree sign
    text = text.replace("°", "")

    # unwrap \text{...}
    match = re.match(r"^\\text\{(?P<x>.+?)\}$", text)
    if match:
        text = match.group("x")

    # remove delimiters \( \) and \[ \]
    text = re.sub(r"\\\(|\\\)|\\\[|\\\]", "", text)

    # apply latex fixes
    for pat, rep in LATEX_FIXES:
        text = re.sub(pat, rep, text)

    def convert_superscripts(s, base=None):
        converted = "".join(
            SUPERSCRIPT_MAP[ch] if ch in SUPERSCRIPT_MAP else ch for ch in s
        )
        if base is None:
            return converted
        return f"{base}**{converted}"

    # transform superscript characters following an alphanumeric or closing bracket
    text = re.sub(
        r"([0-9A-Za-z\)\]\}])([⁰¹²³⁴⁵⁶⁷⁸⁹⁺⁻]+)",
        lambda m: convert_superscripts(m.group(2), base=m.group(1)),
        text,
    )
    text = convert_superscripts(text)

    # replace percent and dollar signs, then remove remaining percent signs
    text = text.replace("\\%", "%").replace("$", "").replace("%", "")

    # sqrt conversions
    text = re.sub(
        r"\\sqrt\s*\{([^}]*)\}", lambda match: f"sqrt({match.group(1)})", text
    )
    text = re.sub(
        r"\\sqrt\s+([^\\\s{}]+)", lambda match: f"sqrt({match.group(1)})", text
    )

    # frac conversions
    text = re.sub(
        r"\\frac\s*\{([^{}]+)\}\s*\{([^{}]+)\}",
        lambda match: f"({match.group(1)})/({match.group(2)})",
        text,
    )
    text = re.sub(
        r"\\frac\s+([^\s{}]+)\s+([^\s{}]+)",
        lambda match: f"({match.group(1)})/({match.group(2)})",
        text,
    )

    # caret to python exponent
    text = text.replace("^", "**")

    # convert cases like "5 1/2" -> "5+1/2"
    text = re.sub(r"(?<=\d)\s+(\d+/\d+)", lambda match: "+" + match.group(1), text)

    # remove thousands separators like "1,000"
    text = re.sub(r"(?<=\d),(?=\d\d\d(\D|$))", "", text)

    return text.replace("{", "").replace("}", "").strip().lower()


print(normalize_text(extract_final_candidate(model_answer)))
from sympy.parsing import sympy_parser as spp
from sympy.core.sympify import SympifyError
from sympy.polys.polyerrors import PolynomialError
from tokenize import TokenError


def sympy_parser(expr):
    """Parse a math expression string into a SymPy expression.

    Returns None on parse errors or if the expression is empty/too long.
    """
    if expr is None or len(expr) > 2000:  # A
        return None

    try:
        return spp.parse_expr(
            expr,
            transformations=(
                *spp.standard_transformations,  # B
                # C
                spp.implicit_multiplication_application,
            ),
            evaluate=True,  # D
        )
    except (
        SympifyError,
        SyntaxError,
        TypeError,
        AttributeError,
        IndexError,
        TokenError,
        ValueError,
        PolynomialError,
    ):
        return None


print(sympy_parser(normalize_text(extract_final_candidate(model_answer))))
print(sympy_parser("28/6"))

from sympy import simplify


def equality_check(expr_gtruth, expr_pred):
    if expr_gtruth == expr_pred:
        return True

    gtruth, pred = sympy_parser(expr_gtruth), sympy_parser(expr_pred)

    if gtruth is not None and pred is not None:
        try:
            return simplify(gtruth - pred) == 0
        except (SympifyError, TypeError):
            pass

    return False


print(equality_check(normalize_text("13/4."), normalize_text(r"(13)/(4)")))

print(equality_check(normalize_text("0.5"), normalize_text(r"(1)/(2)")))


def split_into_parts(text):
    result = [text]
    if text:
        if (
            # A
            len(text) >= 2
            and text[0] in "(["
            and text[-1] in ")]"
            and "," in text[1:-1]
        ):
            items = [p.strip() for p in text[1:-1].split(",")]  # B
            if all(items):
                result = items
            else:  # C
                result = []
    return result


print(split_into_parts(normalize_text(r"(14/3, 2/3)")))


def grade_answer(pred_text, gt_text):
    result = False
    if pred_text is not None and gt_text is not None:
        gt_parts = split_into_parts(normalize_text(gt_text))
        pred_parts = split_into_parts(normalize_text(pred_text))

        if (gt_parts and pred_parts and len(gt_parts) == len(pred_parts)):
            result = all(
                equality_check(gt, pred)
                for gt, pred in zip(gt_parts, pred_parts)
            )
    return result


print(grade_answer("14/3", r"\\frac{14}{3}"))


import json
import requests


def load_math500_test(local_path="math500_test.json", save_copy=True):
    local_path = Path(local_path)
    url = (
        "https://raw.githubusercontent.com/rasbt/reasoning-from-scratch/"
        "main/ch03/01_main-chapter-code/math500_test.json"
    )
    if local_path.exists():
        with local_path.open("r", encoding="utf-8") as f:
            data = json.load(f)
    else:
        r = requests.get(url, timeout=30)
        r.raise_for_status()
        data = r.json()
        if save_copy:  # Saves a local copy
            with local_path.open("w", encoding="utf-8") as f:
                json.dump(data, f, indent=2)
    return data


math_data = load_math500_test()
print("Number of entries:", len(math_data))

from pprint import pprint
pprint(math_data[0])


def render_prompt(prompt):
    template = (
        "You are a helpful math assistant.\n"
        "Answer the question and write the final result on a new line as:\n"
        "\\boxed{ANSWER}\n\n"
        f"Question:\n{prompt}\n\nAnswer:"
    )
    return template


prompt = (
    r"If $a+b=3$ and $ab=\tfrac{13}{6}$, "
    r"what is the value of $a^2+b^2$?"
)

prompt_fmt = render_prompt(prompt)
print(prompt_fmt)

generated_text = generate_text_stream_concat(
    model, tokenizer, prompt_fmt, device,
    max_new_tokens=2048,
    verbose=True
)


def mini_eval_demo(model, tokenizer, device):
    ex = {  # A
        "problem": "Compute 1/2 + 1/6.",
        "answer": "2/3"
    }
    prompt = render_prompt(ex["problem"])
    gen_text = generate_text_stream_concat(
        model, tokenizer, prompt, device,
        max_new_tokens=64
    )
    pred_answer = extract_final_candidate(gen_text)
    is_correct = grade_answer(
        pred_answer, ex["answer"]
    )

    print(f"Device: {device}")
    print(f"Prediction: {pred_answer}")
    print(f"Ground truth: {ex['answer']}")
    print(f"Correct: {is_correct}")


mini_eval_demo(model, tokenizer, device)


import time


def eta_progress_message(
    processed,
    total,
    start_time,
    show_eta=False,
    label="Progress",
):
    progress = f"{label}: {processed}/{total}"
    pad_width = len(f"{label}: {total}/{total} | ETA: 00h 00m 00s")
    if not show_eta or processed <= 0:
        return progress.ljust(pad_width)
    elapsed = time.time() - start_time
    if elapsed <= 0:
        return progress.ljust(pad_width)
    remaining = max(total - processed, 0)
    if processed:
        avg_time = elapsed / processed
        eta_seconds = avg_time * remaining
    else:
        eta_seconds = 0
    eta_seconds = max(int(round(eta_seconds)), 0)
    minutes, rem_seconds = divmod(eta_seconds, 60)
    hours, minutes = divmod(minutes, 60)
    if hours:
        eta = f"{hours}h {minutes:02d}m {rem_seconds:02d}s"
    elif minutes:
        eta = f"{minutes:02d}m {rem_seconds:02d}s"
    else:
        eta = f"{rem_seconds:02d}s"
    message = f"{progress} | ETA: {eta}"
    return message.ljust(pad_width)


def evaluate_math500_stream(
    model,
    tokenizer,
    device,
    math_data,
    out_path=None,
    max_new_tokens=512,
    verbose=False,
):
    if out_path is None:
        dev_name = str(device).replace(":", "-")
        out_path = Path(f"math500-{dev_name}.jsonl")
    num_examples = len(math_data)
    num_correct = 0
    start_time = time.time()
    with open(out_path, "w", encoding="utf-8") as f:
        for i, row in enumerate(math_data, start=1):
            prompt = render_prompt(row["problem"])
            gen_text = generate_text_stream_concat(
                model, tokenizer, prompt, device,
                max_new_tokens=max_new_tokens,
                verbose=verbose,
            )
            extracted = extract_final_candidate(gen_text)
            is_correct = grade_answer(
                extracted, row["answer"]
            )
            num_correct += int(is_correct)

            record = {
                "index": i,
                "problem": row["problem"],
                "gtruth_answer": row["answer"],
                "generated_text": gen_text,
                "extracted": extracted,
                "correct": bool(is_correct),
            }
            f.write(json.dumps(record, ensure_ascii=False) + "\n")
            progress_msg = eta_progress_message(
                processed=i,
                total=num_examples,
                start_time=start_time,
                show_eta=True,
                label="MATH-500",
            )
            print(progress_msg, end="\r", flush=True)
            if verbose:
                # I
                print(
                    f"\n\n{'='*50}\n{progress_msg}\n"
                    f"{'='*50}\nExtracted: {extracted}\n"
                    f"Expected: {row['answer']}\n"
                    f"Correct so far: {num_correct}\n{'-'*50}"
                )
    seconds_elapsed = time.time() - start_time
    acc = num_correct / num_examples if num_examples else 0.0
    print(f"\nAccuracy: {acc*100:.1f}% ({num_correct}/{num_examples})")
    print(f"Total time: {seconds_elapsed/60:.1f} min")
    print(f"Logs written to: {out_path}")
    return num_correct, num_examples, acc

print("Model:", WHICH_MODEL)
num_correct, num_examples, acc = evaluate_math500_stream(
    model, tokenizer, device,
    math_data=math_data[:10],
    max_new_tokens=2048,
    verbose=True
)
