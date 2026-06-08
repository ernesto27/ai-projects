import json
from pathlib import Path
from pprint import pprint

import requests
import torch

from utils import (
    generate_text_stream_concat_flex,
    generate_text_temp_stream_cache,
    get_device,
    load_model_and_tokenizer,
    render_prompt,
    sample_response,
)

device = get_device()
device = torch.device("cpu")

model, tokenizer = load_model_and_tokenizer(
    which_model="base", device=device, use_compile=False
)

# raw_prompt = "Half the value of $3x-9$ is $x+37$. What is the value of $x$?"
# prompt = render_prompt(raw_prompt)
# torch.manual_seed(0)

# response = generate_text_stream_concat_flex(
#     model,
#     tokenizer,
#     prompt,
#     device,
#     max_new_tokens=2048,
#     verbose=True,
#     generate_func=generate_text_temp_stream_cache,
#     temperature=0.9,
#     top_p=0.9,
# )


def load_math_train(local_path="math_train.json", save_copy=True):
    local_path = Path(local_path)
    url = (
        "https://raw.githubusercontent.com/rasbt/"
        "math_full_minus_math500/refs/heads/main/"
        "math_full_minus_math500.json"
    )
    backup_url = (  # A
        "https://f001.backblazeb2.com/file/reasoning-from-scratch/"
        "MATH/math_full_minus_math500.json"
    )

    if local_path.exists():  # B
        with local_path.open("r", encoding="utf-8") as f:
            data = json.load(f)
    else:
        try:
            r = requests.get(url, timeout=30)
            r.raise_for_status()
        except requests.RequestException:
            print("Using backup URL.")
            r = requests.get(backup_url, timeout=30)
            r.raise_for_status()
        data = r.json()

    if save_copy:
        with local_path.open("w", encoding="utf-8") as f:  # C
            json.dump(data, f, indent=2)

    return data


# math_train = load_math_train()
# pprint(math_train[4])

raw_prompt = "Half the value of $3x-9$ is $x+37$. What is the value of $x$?"
prompt = render_prompt(raw_prompt)

token_ids, prompt_len, answer_text = sample_response(
    model=model,
    tokenizer=tokenizer,
    prompt=prompt,
    device=device,
    max_new_tokens=512,
    temperature=0.9,
    top_p=0.9,
)

print(answer_text)

rollouts = [
    r"\boxed{83}",
    r"The correct answer is \boxed{83}",
    r"The final answer is 83",
    r"We get \boxed{38}",
]
