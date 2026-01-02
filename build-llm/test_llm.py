import tiktoken
import torch
import torch.nn as nn
from llm import DummyGPTModel, FeedForward,LayerNorm, ExampleDeepNeuralNetwork, TransformerBlock,  GPT_CONFIG_124M, GPTModel

tokenizer = tiktoken.get_encoding("gpt2")

batch = []
txt1 = "Every effort moves you"
txt2 = "Every day holds a"

batch.append(torch.tensor(tokenizer.encode(txt1)))
batch.append(torch.tensor(tokenizer.encode(txt2)))

batch = torch.stack(batch, dim=0)

torch.manual_seed(123) 
model = GPTModel(GPT_CONFIG_124M)

out = model(batch)
print("Input batch: \n", batch)
print("Output shape:\n", out.shape)
print(out)

total_params = sum(p.numel() for p in model.parameters())
total_params_gpt2 = (total_params - sum(p.numel() for p in model.out_head.parameters()))
print(f"Number of trainable parameters " f"considering weight tying: {total_params_gpt2:,}")


total_size_bytes = total_params * 4
total_size_mb = total_size_bytes / (1024 * 1024)
print(f"Total size of the model: {total_size_mb:.2f} MB")




def generate_text_simple(model, idx, max_new_tokens, context_size):
    for _ in range(max_new_tokens):
        idx_cond = idx[:, -context_size:]
        with torch.no_grad():
            logits = model(idx_cond)
        
        logits = logits[:, -1, :] 
        probas = torch.softmax(logits, dim=-1)
        idx_next = torch.argmax(probas, dim=-1, keepdim=True)
        idx = torch.cat((idx, idx_next), dim=1)
    return idx

start_context = "hello, i am"
encoded = tokenizer.encode(start_context) 
print("Encoded input:", encoded)
encoded_tensor = torch.tensor(encoded).unsqueeze(0)
print("Encoded tensor shape:", encoded_tensor.shape)

model.eval()
out = generate_text_simple(
    model=model,
    idx=encoded_tensor,
    max_new_tokens=6, 
    context_size=GPT_CONFIG_124M["context_length"]
)

print("output:", out)
print("output lenght:", len(out[0]))

decoded_text = tokenizer.decode(out.squeeze(0).tolist())
print(decoded_text)

# torch.manual_seed(123)
# x = torch.randn(2, 4, 768)
# block = TransformerBlock(GPT_CONFIG_124M)
# out = block(x)
# print("input shape:", x.shape)
# print("Transformer block output shape:", out.shape)



# layer_sizes = [3,3,3,3,3,1]
# sample_output = torch.tensor([[1., 0., -1.]])
# torch.manual_seed(123)
# model_without_shortcut = ExampleDeepNeuralNetwork(layer_sizes, use_shorcut=True)

# def print_gradients(model, x):
#     output = model(x)
#     target = torch.tensor([[0.]])
#     loss = nn.MSELoss()
#     loss = loss(output, target)
#     loss.backward()

#     for name, param in model.named_parameters():
#         if 'weight' in name:
#             print(f"{name} has gradient mean of {param.grad.abs().mean().item()}")

# print_gradients(model_without_shortcut, sample_output)

# tokenizer = tiktoken.get_encoding("gpt2")
# batch = []
# txt1 = "Every effort moves you"
# txt2 = "Every day holds a"

# batch.append(torch.tensor(tokenizer.encode(txt1)))
# batch.append(torch.tensor(tokenizer.encode(txt2)))

# batch = torch.stack(batch, dim=0)
# print(batch)

# torch.manual_seed(123)
# model = DummyGPTModel(GPT_CONFIG_124M)
# logits = model(batch)
# print("output shape:\n", logits.shape)  # Expected: (2, seq
# print(logits)

# batch_example = torch.randn(2, 5) 
# layer = nn.Sequential(nn.Linear(5,6), nn.ReLU()) 
# out = layer(batch_example)
# print(out)

# mean = out.mean(dim=-1, keepdim=True)
# var = out.var(dim=-1, keepdim=True)
# print("Mean:\n", mean)
# print("Variance:\n", var)


# ln = LayerNorm(emb_dim=5)
# out_ln = ln(batch_example)
# mean = out_ln.mean(dim=-1, keepdim=True)
# var = out_ln.var(dim=-1, unbiased=False, keepdim=True)
# print("Mean:\n", mean)
# print("Variance:\n", var)

# ffn = FeedForward(GPT_CONFIG_124M)
# x = torch.randn(2, 3, 768) 
# out = ffn(x) 
# print(out.shape)