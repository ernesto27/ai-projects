import json

import torch
import matplotlib.pyplot as plt
from reasoning_from_scratch.qwen3 import KVCache
from utils import (
    extract_final_candidate,
    get_device,
    generate_text_basic_stream,
    load_model_and_tokenizer,
    render_prompt,
    generate_text_stream_concat_flex,
)


device = get_device()
device = torch.device("cpu")

model, tokenizer = load_model_and_tokenizer(
    which_model="base", device=device, use_compile=False
)

raw_prompt = "Half the value of $3x-9$ is $x+37$. What is the value of $x$?"

prompt = render_prompt(raw_prompt)
print(prompt)


# response = generate_text_stream_concat_flex(
#     model, tokenizer, prompt, device,
#     max_new_tokens=2048, verbose=True,
#     generate_func=generate_text_basic_stream
# )

# prompt_cot = prompt + " \n\nExplain step by step."

# response_cot = generate_text_stream_concat_flex(
#     model, tokenizer, prompt_cot, device,
#     max_new_tokens=2048, verbose=True
# )


ex_prompt = "The capital of Germany is"
response = generate_text_stream_concat_flex(
    model, tokenizer, ex_prompt, device, max_new_tokens=1, verbose=True
)


input_token_ids = torch.tensor(tokenizer.encode(ex_prompt), device=device).unsqueeze(0)
print(input_token_ids)

with torch.inference_mode():
    next_token_logits = model(input_token_ids)[:, -1]
print(next_token_logits.shape)


max_token_id = torch.argmax(next_token_logits)
print(f"Token ID: {max_token_id}")
print(f"Decoded token: '{tokenizer.decode([max_token_id])}'")


def plot_scores_bar(
    next_token_logits, start=19_800, end=19_900, arrow=True, ylabel="Logit value"
):
    x = torch.arange(start, end)
    logits_section = next_token_logits[0, start:end].float().cpu()
    plt.bar(x, logits_section)
    plt.xlabel("Vocabulary index")
    plt.ylabel(ylabel)
    if arrow:
        max_idx = torch.argmax(logits_section)
        plt.annotate(
            "Berlin",
            xy=(x[max_idx], logits_section[max_idx]),
            xytext=(x[max_idx] - 25, logits_section[max_idx] - 2),
            arrowprops={"facecolor": "black", "arrowstyle": "->", "lw": 1.5},
            fontsize=10,
        )
    plt.grid(alpha=0.3)
    plt.tight_layout()
    plt.show()


# plot_scores_bar(next_token_logits)


def scale_logits_by_temperature(logits, temperature):
    if temperature < 0:
        raise ValueError("Temeperature must be positive")
    return logits / temperature


def plot_logits_with_temperature(
    next_token_logits,
    start=19_800,
    end=19_900,
    temps=(0.5, 5.0),
):
    x = torch.arange(start, end)
    logits_orig = next_token_logits[0, start:end].float().cpu()
    logits_scaled = [scale_logits_by_temperature(logits_orig, T) for T in temps]
    plt.plot(x, logits_orig, label="Original logits", lw=2)
    plt.plot(x, logits_scaled[0], label=f"T={temps[0]} (sharper)", ls="--", lw=1)
    plt.plot(x, logits_scaled[1], label=f"T={temps[1]} (flatter)", ls=":", lw=3)

    max_idx = torch.argmax(logits_orig)
    plt.annotate(
        "Berlin",
        xy=(x[max_idx], logits_orig[max_idx]),
        xytext=(x[max_idx] - 25, logits_orig[max_idx] + 2),
        arrowprops={"facecolor": "black", "arrowstyle": "->", "lw": 1.5},
        fontsize=12,
    )
    plt.xlabel("Vocabulary index")
    plt.ylabel("Logit value")
    plt.legend()
    plt.grid(alpha=0.3)
    plt.tight_layout()
    plt.show()


# plot_logits_with_temperature(next_token_logits, temps=(0.5, 5.0))

# rescaled_logits = scale_logits_by_temperature(next_token_logits, 5.0)
rescaled_logits = scale_logits_by_temperature(next_token_logits, 0.35)

next_token_probas = torch.softmax(rescaled_logits, dim=-1)

print(torch.sum(next_token_probas))
print("Token ID 19,846 probability:", next_token_probas[:, 19846])

torch.manual_seed(123)
print("sampled token: ", torch.multinomial(next_token_probas.cpu(), num_samples=1))


def count_samples(probas, num_samples=1000, threshold=1, tokenizer=None):
    samples = torch.multinomial(probas.cpu(), num_samples=num_samples, replacement=True)
    counts = torch.bincount(samples.squeeze(0), minlength=1)

    for i, c in enumerate(counts):
        if c > threshold:
            if tokenizer is None:
                print(f"Vocab index {i}: {c.item()}x")
            else:
                print(f"'{tokenizer.decode([i])}': {c.item()}x")


count_samples(next_token_probas, tokenizer=tokenizer)
print(next_token_probas[0, 19_846])


@torch.inference_mode()
def generate_text_temp_stream_cache(
    model, token_ids, max_new_tokens, eos_token_id=None, temperature=0.0, top_p=None
):
    model.eval()
    cache = KVCache(n_layers=model.cfg["n_layers"])
    model.reset_kv_cache()

    out = model(token_ids, cache=cache)[:, -1]

    for _ in range(max_new_tokens):
        orig_device = token_ids.device

        if temperature is None or temperature == 1.0:
            next_token = torch.argmax(out, dim=-1, keepdim=True)
        else:
            logits = scale_logits_by_temperature(out, temperature)
            probas = torch.softmax(logits, dim=-1)

            probas = top_p_filter(probas, top_p)

            next_token = torch.multinomial(probas.cpu(), num_samples=1)
            next_token = next_token.to(orig_device)

        if eos_token_id is not None and torch.all(next_token == eos_token_id):
            break

        yield next_token
        out = model(next_token, cache=cache)[:, -1]


# torch.manual_seed(123)
# response = generate_text_stream_concat_flex(
#     model, tokenizer, prompt, device ,
#     max_new_tokens=2048, verbose=True,
#     generate_func=generate_text_temp_stream_cache,
#     temperate=1.1
# )

# print(response)

toy_logits = torch.tensor(  # A
    [-0.7, -3.0, 0.1, -1.2, 2.0, -1.0, -0.5, -2.0, 0.3, 1.5]
)

toy_logits_scaled = scale_logits_by_temperature(toy_logits, 1.0)
toy_probas = torch.softmax(toy_logits_scaled, dim=-1)

# plt.bar(
#     torch.arange(len(toy_probas)), toy_probas, alpha=0.5
# )

# plt.ylim([0, 1])
# plt.xlabel("Vocabulary index")
# plt.ylabel("Probability")
# plt.show()

sorted_probas, sorted_idx = torch.sort(toy_probas, descending=True)
cumsum = torch.cumsum(sorted_probas, dim=-1)

# plt.bar(
#     torch.arange(len(sorted_probas)), sorted_probas,
#     alpha=0.5
# )
# plt.step(
#     torch.arange(len(cumsum)), cumsum,
#     where="mid", color="C1", label="Cumulative sum"
# )
# plt.ylim([0, 1])
# plt.xlabel("Token rank (sorted by probability)")
# plt.ylabel("Probability")
# plt.show()

top_p = 0.8
keep_mask = cumsum <= top_p
n_kept = torch.sum(keep_mask).item()
print("cumulative sum:", cumsum)
print("tokens kept:", n_kept)

keep_mask = (cumsum - sorted_probas) < top_p
n_kept = keep_mask.sum().item()
print("Tokens kept:", n_kept)

kept_sorted = torch.where(keep_mask, sorted_probas, torch.zeros_like(sorted_probas))
filtered = torch.zeros_like(toy_probas).scatter(0, sorted_idx, kept_sorted)
print(filtered)

denom = torch.sum(filtered).clamp_min(1e-12)
renormalized = filtered / denom
print(renormalized)


def top_p_filter(probas, top_p):
    if top_p is None or top_p >= 1.0:
        return probas

    sorted_probas, sorted_idx = torch.sort(probas, dim=1, descending=True)
    cumprobas = torch.cumsum(sorted_probas, dim=1)

    prefix = cumprobas - sorted_probas
    keep = prefix < top_p
    keep[:, 0] = True

    keep_sorted = torch.where(keep, sorted_probas, torch.zeros_like(sorted_probas))

    filtered = torch.zeros_like(probas).scatter(1, sorted_idx, keep_sorted)

    denom = torch.sum(filtered, dim=1, keepdim=True).clamp_min(1e-12)
    return filtered / denom


with torch.inference_mode():
    next_token_logits = model(input_token_ids)[:, -1]
print(next_token_logits.shape)

probas_lowT = torch.softmax(
    scale_logits_by_temperature(next_token_logits, 0.35), dim=-1
)
probas_lowT_filtered = top_p_filter(probas_lowT, top_p=0.8)
count_samples(probas_lowT_filtered, threshold=1, tokenizer=tokenizer)


torch.manual_seed(123)
response = generate_text_stream_concat_flex(
    model,
    tokenizer,
    prompt,
    device,
    max_new_tokens=2048,
    verbose=True,
    generate_func=generate_text_temp_stream_cache,
    temperature=0.5,
    top_p=0.8,
)

from collections import Counter


def self_consistency_vote(
    model,
    tokenizer,
    prompt,
    device,
    num_samples=10,
    temperature=0.8,
    top_p=0.9,
    max_new_tokens=2048,
    show_progress=True,
    show_long_answer=False,
    seed=None,
):
    full_answers, short_answers = [], []
    for i in range(num_samples):  # A
        if seed is not None:
            torch.manual_seed(seed + i + 1)
        answer = generate_text_stream_concat_flex(
            model=model,
            tokenizer=tokenizer,
            prompt=prompt,
            device=device,
            max_new_tokens=max_new_tokens,
            verbose=show_long_answer,
            generate_func=generate_text_temp_stream_cache,
            temperature=temperature,
            top_p=top_p,
        )
        short = extract_final_candidate(
            # B
            answer,
            fallback="number_then_full",
            # B
        )
        # B
        full_answers.append(answer)
        short_answers.append(short)
        if show_progress:
            print(f"[Sample {i + 1}/{num_samples}] → {short!r}")
    counts = Counter(short_answers)
    groups = {s: [] for s in counts}
    for idx, s in enumerate(short_answers):
        groups[s].append(idx)
    mc = counts.most_common()  # C
    if not mc:
        majority_winners, final_answer = [], None
    else:
        top_freq = mc[0][1]
        majority_winners = [s for s, f in mc if f == top_freq]
        final_answer = mc[0][0] if len(majority_winners) == 1 else None

    return {
        "full_answers": full_answers,
        "short_answers": short_answers,
        "counts": dict(counts),
        "groups": groups,
        "majority_winners": majority_winners,
        "final_answer": final_answer,
    }


results = self_consistency_vote(
    model,
    tokenizer,
    prompt,
    device=device,
    num_samples=5,
    temperature=0.8,N
    top_p=0.9,
    max_new_tokens=2048,
    seed=123,
    show_progress=True,
)

print(json.dumps(results, indent=2, ensure_ascii=False, default=str))
