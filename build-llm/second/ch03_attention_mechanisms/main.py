import torch
inputs = torch.tensor(
    [[0.43, 0.15, 0.89], # Your
    [0.55, 0.87, 0.66], # journey
    [0.57, 0.85, 0.64], # starts
    [0.22, 0.58, 0.33], # with
    [0.77, 0.25, 0.10], # one
    [0.05, 0.80, 0.55]] # step
)

query = inputs[1]
attn_scores_2 = torch.empty(inputs.shape[0]) 
for i, x_i in enumerate(inputs):
    attn_scores_2[i] = torch.dot(x_i, query)

print("Attention scores for the second token (journey):\n", attn_scores_2)
for i, score in enumerate(attn_scores_2):
    print(f"attn_scores_2[{i}] = {score.item():.4f}")

attn_weights_2_tmp = attn_scores_2 / attn_scores_2.sum()
print("Attention weights:", attn_weights_2_tmp)
print("Sum:", attn_weights_2_tmp.sum())

def sotfmax_naive(x):
    return torch.exp(x) / torch.exp(x).sum(dim=0)

attn_weights_2_naive = sotfmax_naive(attn_scores_2)
print("Attention weights (naive softmax):", attn_weights_2_naive)
print("Sum (naive softmax):", attn_weights_2_naive.sum())

attn_weights_2 = torch.softmax(attn_scores_2, dim=0)
print("Attention weights:", attn_weights_2)
print("Sum:", attn_weights_2.sum())

query = inputs[1]
context_vec_2 = torch.zeros(query.shape)
for i, x_i in enumerate(inputs):
    context_vec_2 += attn_weights_2[i] * x_i        

print("Context vector for the second token (journey):\n", context_vec_2)

attn_scores = torch.empty(6,6)
for i, x_i in enumerate(inputs):
    for j, x_j in enumerate(inputs):
        attn_scores[i, j] = torch.dot(x_i, x_j)

# print("Attention scores matrix:\n", attn_scores)

# attn_scores = inputs @ inputs.T
# print("Attention scores matrix (using matrix multiplication):\n", attn_scores)

# attn_weights = torch.softmax(attn_scores, dim=1)
# print("Attention weights matrix:\n", attn_weights)

# row_2_sum = sum([0.1385, 0.2379, 0.2333, 0.1240, 0.1082, 0.1581])
# print("Row 2 sum:", row_2_sum)
# print("All row sums:", attn_weights.sum(dim=-1))

# all_context_vecs = attn_weights @ inputs
# print(all_context_vecs)

# 3.4s

x_2 = inputs[1]
d_in = inputs.shape[1]
d_out = 2

torch.manual_seed(123)
W_query = torch.nn.Parameter(torch.randn(d_in, d_out), requires_grad=False)
W_key = torch.nn.Parameter(torch.randn(d_in, d_out), requires_grad=False)
W_value = torch.nn.Parameter(torch.randn(d_in, d_out), requires_grad=False)

query_2 = x_2 @ W_query
key_2 = x_2 @ W_key
value_2 = x_2 @ W_value
print(query_2)


keys = inputs @ W_key
values = inputs @ W_value
print("keys.shape:", keys.shape)
print("values.shape:", values.shape)

keys_2 = keys[1]
attn_scores_22 = query_2.dot(keys_2)
print(attn_scores_22)

attn_scores_2 = query_2 @ keys.T
print(attn_scores_2)

d_k = keys.shape[-1]
attn_weights_2 = torch.softmax(attn_scores_2 / d_k**0.5, dim=-1)
print(attn_weights_2)

context_vec_2 = attn_weights_2 @ values
print(context_vec_2)


