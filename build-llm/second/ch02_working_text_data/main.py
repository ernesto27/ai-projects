import urllib.request
import os

# url = ("https://raw.githubusercontent.com/rasbt/LLMs-from-scratch/main/ch02/01_main-chapter-code/the-verdict.txt")
# file_path = "the-verdict.txt"
# urllib.request.urlretrieve(url, file_path)

with open(os.path.join(os.path.dirname(__file__), "the-verdict.txt"), "r", encoding="utf-8") as f:
    raw_text = f.read()

# print("Total characters in the text:", len(raw_text))
# print(raw_text[:99])


import re 
# text = "Hello, world. Is this-- a test?"
# result = re.split(r'([,.:;?_!"()\']|--|\s)', text)
# result = [item for item in result if item.strip()]
# print(result)

preprocessed = re.split(r'([,.:;?_!"()\']|--|\s)', raw_text)
preprocessed = [item.strip() for item in preprocessed if item.strip()]
print(len(preprocessed))
print(preprocessed[:30])


all_words = sorted(set(preprocessed))
vocab_size = len(all_words)
print("Vocabulary size:", vocab_size)


vocab = {token:integer for integer, token in enumerate(all_words)}
for i, item in enumerate(vocab.items()):
    print(item)
    if i >= 50:
        break


class SimpleTokenizerV1:
    def __init__(self, vocab):
        self.str_to_int = vocab
        self.int_to_str = {integer: token for token, integer in vocab.items()}

    def encode(self, text):
        preprocessed = re.split(r'([,.:;?_!"()\']|--|\s)', text)
        preprocessed = [item.strip() for item in preprocessed if item.strip()]
        preprocessed = [item if item in self.str_to_int else "<|unk|>" for item in preprocessed]
        ids = [self.str_to_int[s] for s in preprocessed]
        return ids

    def decode(self, ids):
        text = " ".join([self.int_to_str[i] for i in ids])
        text = re.sub(r'\s+([,.:;?!"()\'])', r'\1', text)
       
        return text


tokenizer = SimpleTokenizerV1(vocab)

# text = """"It's the last he painted, you know,"
# Mrs. Gisburn said with pardonable pride."""
# ids = tokenizer.encode(text)
# print(ids)
# print(tokenizer.decode(ids))

# text = "Hello, do you like tea?"
# print(tokenizer.encode(text))

all_tokens = sorted(list(set(preprocessed)))
all_tokens.extend(["<|endoftext|>", "<|unk|>"])
vocab = {token:integer for integer, token in enumerate(all_tokens)}
# print(len(vocab.items()))

# for i, item in enumerate(list(vocab.items())[-5:]):
#     print(item)


# text1 = "Hello, do you like tea?"
# text2 = "In the sunlit terraces of the palace."
# text = " <|endoftext|> ".join((text1, text2))
# print(text)

# tokenizer = SimpleTokenizerV1(vocab)
# print(tokenizer.encode(text))

# print(tokenizer.decode(tokenizer.encode(text)))

from importlib.metadata import version
import tiktoken
print("tiktoken version:", version("tiktoken"))
tokenizer = tiktoken.get_encoding("gpt2")
text = (
"Hello, do you like tea? <|endoftext|> In the sunlit terraces"
"of someunknownPlace."
)
integers = tokenizer.encode(text, allowed_special={"<|endoftext|>"})
print(integers)

strings = tokenizer.decode(integers)
print(strings)


with open(os.path.join(os.path.dirname(__file__), "the-verdict.txt"), "r", encoding="utf-8") as f:
    raw_text = f.read()

enc_text = tokenizer.encode(raw_text)
print(len(enc_text))

enc_sample = enc_text[50:]
context_size = 4 
x = enc_sample[:context_size]
y = enc_sample[1:context_size+1]
print(f"x: {x}")
print(f"y:      {y}")

for i in range(1, context_size+1): 
    context = enc_sample[:i]
    target = enc_sample[i]
    print(context, "----->", target)
    print(tokenizer.decode(context), "----->", tokenizer.decode([target]))


import torch
from  torch.utils.data import Dataset, DataLoader

class GPTDatasetV1(Dataset):
    def __init__(self, txt, tokenizer, max_length, stride):
        self.input_ids = []
        self.target_ids = []

        token_ids = tokenizer.encode(txt) 

        for i in range(0, len(token_ids) - max_length, stride): 
            input_chunk = token_ids[i:i + max_length]
            target_chunk = token_ids[i + 1:i + max_length + 1]
            self.input_ids.append(torch.tensor(input_chunk))
            self.target_ids.append(torch.tensor(target_chunk))

    def __len__(self):
        return len(self.input_ids)

    def __getitem__(self, idx):
        return self.input_ids[idx], self.target_ids[idx]


def create_dataloader(txt, batch_size=4, max_length=256, stride=128, shuffle=True, drop_last=True, num_workers=0): 
    tokenizer = tiktoken.get_encoding("gpt2") 
    dataset = GPTDatasetV1(txt, tokenizer, max_length, stride) 
    dataloader = DataLoader( 
        dataset, 
        batch_size=batch_size, 
        shuffle=shuffle, 
        drop_last=drop_last,
        num_workers=num_workers
    )

    return dataloader

with open(os.path.join(os.path.dirname(__file__), "the-verdict.txt"), "r", encoding="utf-8") as f:
    raw_text = f.read()

dataloader = create_dataloader(raw_text, batch_size=1, max_length=4, stride=1, shuffle=False)
data_iter = iter(dataloader)
first_batch = next(data_iter) 
print(first_batch)

second_batch = next(data_iter)
print(second_batch)

input_ids = torch.tensor([2, 3, 5, 1])
vocab_size = 6 
ouput_dim = 3
torch.manual_seed(123)
embedding_layer = torch.nn.Embedding(vocab_size, ouput_dim)
# print(embedding_layer.weight)
# # print(embedding_layer(torch.tensor([3])))
# print("---------------")
# print(embedding_layer(input_ids))


vocab_size = 50257
output_dim = 256 

token_embedding_layer = torch.nn.Embedding(vocab_size, output_dim)

max_length = 4
dataloader = create_dataloader(raw_text, batch_size=8, max_length=max_length, stride=max_length, shuffle=False)
data_iter = iter(dataloader)
inputs, targets = next(data_iter)
print("Token IDs:\n", inputs)
print("\nInputs shape:\\n", inputs.shape)

token_embeddings = token_embedding_layer(inputs)
print("\nToken Embeddings shape:\n", token_embeddings.shape)


context_length = max_length
pos_embedding_layer = torch.nn.Embedding(context_length, output_dim)
pos_embeddings = pos_embedding_layer(torch.arange(context_length))
print(pos_embeddings.shape)

input_embeddings = token_embeddings + pos_embeddings 
print("\nInput Embeddings shape:\n", input_embeddings.shape)
