import torch 

print(torch.__version__)
print(torch.cuda.is_available())

tensor0d = torch.tensor(1)

tensor1d = torch.tensor([1, 2, 3])

tensor2d = torch.tensor([[1, 2],[3, 4]])

tensor3d = torch.tensor([[[1, 2], [3, 4]],[[5, 6], [7, 8]]])

print(tensor3d.shape)

print(tensor2d @ tensor2d.T)

import torch.nn.functional as F

y = torch.tensor([1.0])
x1 = torch.tensor([1.1])
w1 = torch.tensor([2.2])
b = torch.tensor([0.0])
z = x1 * w1 + b
a = torch.sigmoid(z)
loss = F.binary_cross_entropy(a, y)