import json

import matplotlib.pyplot as plt
import torch
from reasoning_from_scratch.qwen3 import KVCache

from utils import (
    extract_final_candidate,
    generate_text_basic_stream,
    generate_text_stream_concat_flex,
    generate_text_temp_stream_cache,
    get_device,
    load_model_and_tokenizer,
    render_prompt,
)

device = get_device()
device = torch.device("cpu")

model, tokenizer = load_model_and_tokenizer(
    which_model="base", device=device, use_compile=False
)

raw_prompt = "Half the value of $3x-9$ is $x+37$. What is the value of $x$?"
prompt = render_prompt(raw_prompt)
prompt_cot = prompt + "\n\nExplain step by step."
torch.manual_seed(0)

# response_1 = generate_text_stream_concat_flex(
#     model,
#     tokenizer,
#     prompt_cot,
#     device,
#     max_new_tokens=2048,
#     verbose=True,
#     generate_func=generate_text_temp_stream_cache,
#     temperature=0.9,
#     top_p=0.9,
# )

# torch.manual_seed(3)
# response_2 = generate_text_stream_concat_flex(
#     model,
#     tokenizer,
#     prompt_cot,
#     device,
#     max_new_tokens=2048,
#     verbose=True,
#     generate_func=generate_text_temp_stream_cache,
#     temperature=0.9,
#     top_p=0.9,
# )

# print("Response 1 characters:", len(response_1))
# print("Response 1 tokens:", len(tokenizer.encode(response_1)))
# print("\nResponse 2 characters:", len(response_2))
# print("Response 2 tokens:", len(tokenizer.encode(response_2)))

import math


def heuristic_score(
    answer,
    prompt=None,
    brevity_bunus=500.0,
    boxed_bonus=2.0,
    extract_bonus=1.0,
    fulltext_bonus=0.0,
):
    score = 0.0

    cand = extract_final_candidate(answer, fallback="none")
    if cand:
        score += boxed_bonus
    else:
        cand = extract_final_candidate(answer, fallback="number_only")
        if cand:
            score += extract_bonus
        else:
            cand = extract_final_candidate(answer, fallback="number_then_full")
            if cand:
                score += fulltext_bonus

    score += 1.5 * math.exp(-len(answer) / brevity_bunus)
    return score


# print(round(heuristic_score(response_1), 3))
# print(round(heuristic_score(response_2), 3))


@torch.inference_mode()
def calc_next_token_probas(model, tokenizer, prompt, device, show=True):
    token_ids = torch.tensor(tokenizer.encode(prompt), device=device)
    logits = model(token_ids.unsqueeze(0)).squeeze()
    all_probas = torch.softmax(logits, dim=-1)

    t_idx = torch.arange(0, token_ids.shape[0] - 1, device=device)
    next_ids = token_ids[1:]
    next_token_probas = all_probas[t_idx, next_ids]
    prod_next_token_probas = torch.prod(next_token_probas)

    if show:
        print("Next-token probabilities:", next_token_probas)
        print("Joint probability:", prod_next_token_probas)
    else:
        return next_token_probas, prod_next_token_probas


# torch.set_printoptions(precision=4, sci_mode=True)
# calc_next_token_probas(
#     model, tokenizer, device=device, prompt="The capital of Germany is Berlin"
# )


# calc_next_token_probas(
#     model, tokenizer, device=device, prompt="The capital of Germany is Bridge"
# )


@torch.inference_mode()
def calc_next_token_logprobas(model, tokenizer, prompt, device, show=True):
    token_ids = torch.tensor(tokenizer.encode(prompt), device=device)
    logits = model(token_ids.unsqueeze(0)).squeeze()
    all_probas = torch.log_softmax(logits, dim=-1)

    t_idx = torch.arange(0, token_ids.shape[0] - 1, device=device)
    next_ids = token_ids[1:]
    next_token_probas = all_probas[t_idx, next_ids]
    sum_next_token_probas = torch.sum(next_token_probas)

    if show:
        print("Next-token probabilities:", next_token_probas)
        print("Joint probability:", sum_next_token_probas)
    else:
        return next_token_probas, sum_next_token_probas


calc_next_token_logprobas(
    model, tokenizer, device=device, prompt="The capital of Germany is Berlin"
)
calc_next_token_logprobas(
    model, tokenizer, device=device, prompt="The capital of Germany is Bridge"
)


# example_prompt = "What is the capital of Germany?"
# example_answer = " The capital of Germany is Berlin."
# next_token_logprobas, sum_next_token_logprobas = calc_next_token_logprobas(
#     model, tokenizer, device=device, prompt=example_prompt + example_answer, show=False
# )
# print("11Next-token logprobas:", next_token_logprobas)
# print("22Joint log-probability:", sum_next_token_logprobas)


@torch.inference_mode()
def avg_logprob_answer(model, tokenizer, prompt, answer, device="cpu"):
    prompt_ids = tokenizer.encode(prompt)
    answer_ids = tokenizer.encode(answer)
    full_ids = torch.tensor(prompt_ids + answer_ids, device=device)

    logits = model(full_ids.unsqueeze(0)).squeeze(0)
    logprobs = torch.log_softmax(logits, dim=-1)

    start = len(prompt_ids) - 1
    end = full_ids.shape[0] - 1

    t_idx = torch.arange(start, end, device=device)
    next_tokens = full_ids[start + 1 : end + 1]
    next_token_logps = logprobs[t_idx, next_tokens]

    return float(torch.mean(next_token_logps).item())


# score_1 = avg_logprob_answer(
#     model,
#     tokenizer,
#     prompt="What is the capital of Germany?",
#     answer=" The capital of Germany is Berlin.",
#     device=device,
# )
# print("score1", score_1)


# score_2 = avg_logprob_answer(
#     model,
#     tokenizer,
#     prompt="What is the capital of Germany?",
#     answer=" The capital of Germany is Bridge.",
#     device=device,
# )
# print("score2", score_2)


raw_prompt = "Half the value of $3x-9$ is $x+37$. What is the value of $x$?"
prompt = render_prompt(raw_prompt)

torch.manual_seed(123)
initial_response = generate_text_stream_concat_flex(
    model,
    tokenizer,
    prompt,
    device,
    max_new_tokens=2048,
    verbose=True,
    generate_func=generate_text_temp_stream_cache,
    temperature=0.7,
    top_p=0.9,
)


def make_critique_prompt(raw_prompt, draft):
    return (
        "You are a meticulous reviewer. Identify logical errors, missing "
        "steps, or arithmetic mistakes. If the answer seems correct, "
        "say so briefly. Then propose a concise plan to fix issues.\n\n"
        f"Question:\n{raw_prompt}\n\n"
        f"Draft answer:\n{draft}\n\n"
        "Write a short critique and bullet-point fix plan "
        "(under ~120 words).\n"
        "Critique:"
    )


# critique_prompt = make_critique_prompt(raw_prompt, initial_response)
# critique = generate_text_stream_concat_flex(
#     model,
#     tokenizer,
#     critique_prompt,
#     device,
#     max_new_tokens=2048,
#     verbose=True,
#     generate_func=generate_text_temp_stream_cache,
#     temperature=0.7,
#     top_p=0.9,
# )


def make_refine_prompt(raw_prompt, draft, critique):
    return (
        "Revise the answer using the critique. Keep it concise and "
        "end with a final boxed result: \\boxed{ANSWER}\n\n"
        f"Question:\n{raw_prompt}\n\n"
        f"Previous answer:\n{draft}\n\n"
        f"Critique:\n{critique}\n\n"
        "Revised answer:"
    )


# refine_prompt = make_refine_prompt(raw_prompt, initial_response, critique)
# revise_answer = generate_text_stream_concat_flex(
#     model,
#     tokenizer,
#     refine_prompt,
#     device,
#     max_new_tokens=2048,
#     verbose=True,
#     generate_func=generate_text_temp_stream_cache,
#     temperature=0.7,
#     top_p=0.9,
# )


def self_refinement_loop(
    model,
    tokenizer,
    raw_prompt,
    device,
    iterations=2,
    max_response_tokens=2048,
    max_critique_tokens=256,
    score_fn=None,
    prompt_renderer=render_prompt,
    prompt_suffix="",
    verbose=False,
    temperature=0.7,
    top_p=0.0,
):
    steps = []

    prompt = prompt_renderer(raw_prompt) + prompt_suffix
    current_full = generate_text_stream_concat_flex(
        model=model,
        tokenizer=tokenizer,
        prompt=prompt,
        device=device,
        max_new_tokens=max_response_tokens,
        verbose=False,
        generate_func=generate_text_temp_stream_cache,
        temperature=temperature,
        top_p=top_p,
    )

    current_extracted = extract_final_candidate(
        current_full, fallback="number_then_full"
    )

    if score_fn:
        current_score = score_fn(answer=current_full, prompt=prompt)
    else:
        current_score = 0.0

    for it in range(iterations):
        draft_before_full = current_full
        draft_before_extracted = current_extracted
        score_before = current_score

        critique_prompt = make_critique_prompt(raw_prompt, draft_before_full)
        critique_full = generate_text_stream_concat_flex(
            model=model,
            tokenizer=tokenizer,
            prompt=critique_prompt,
            device=device,
            max_new_tokens=max_critique_tokens,
            verbose=False,
            generate_func=generate_text_temp_stream_cache,
            temperature=temperature,
            top_p=top_p,
        )

        refine_prompt = make_refine_prompt(raw_prompt, draft_before_full, critique_full)

        revised_full = generate_text_stream_concat_flex(
            model=model,
            tokenizer=tokenizer,
            prompt=refine_prompt,
            device=device,
            max_new_tokens=max_response_tokens,
            verbose=False,
            generate_func=generate_text_temp_stream_cache,
            temperature=temperature,
            top_p=top_p,
        )

        revised_extracted = extract_final_candidate(
            revised_full, fallback="number_then_full"
        )
        if score_fn:
            revised_score = score_fn(answer=revised_full, prompt=prompt)
        else:
            revised_score = 0.0

        step = {
            "iteration": it + 1,
            "draft_full": draft_before_full,
            "draft_extracted": draft_before_extracted,
            "critique": critique_full,
            "revised_full": revised_full,
            "revised_extracted": revised_extracted,
            "score_before": score_before,
            "score_after": revised_score,
        }
        steps.append(step)

        if verbose:
            print(
                f"[Refinement {it + 1}/{iterations}]"
                f"\nCurrent: {draft_before_extracted}"
                f"\nRevised: {revised_extracted}"
                f"\nScore before: {score_before:.3f}"
                f"\nScore after: {revised_score:.3f}"
                f"\n{'=' * 25}\n"
            )

        if revised_score >= current_score:
            current_full = revised_full
            current_extracted = revised_extracted
            current_score = revised_score

    return {
        "final_full": current_full,
        "final_extracted": current_extracted,
        "steps": steps,
    }


from functools import partial

avg_logprob_score = partial(
    avg_logprob_answer, model=model, tokenizer=tokenizer, device=device
)


results_logprob = self_refinement_loop(
    model=model,
    tokenizer=tokenizer,
    raw_prompt=raw_prompt,
    device=device,
    iterations=2,
    max_response_tokens=2048,
    max_critique_tokens=256,
    score_fn=avg_logprob_score,
    verbose=True,
    temperature=0.7,
    top_p=0.9,
)
