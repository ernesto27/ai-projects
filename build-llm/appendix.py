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

class NeuralNetwork(torch.nn.Module):
    def __init__(self, num_inputs, num_outputs):
        super().__init__()

        self.layers = torch.nn.Sequential(
            torch.nn.Linear(num_inputs, 30),
            torch.nn.ReLU(), 

            torch.nn.Linear(30,20), 
            torch.nn.ReLU(),

            torch.nn.Linear(20, num_outputs)
        )

        def forward(self, x): 
            logits = self.layers(x)
            return logits

model = NeuralNetwork(50, 3)
print(model)

num_params = sum(p.numel() for p in model.parameters() if p.requires_grad)
print(f"Number of trainable parameters: {num_params}")

X_train = torch.tensor([
[-1.2, 3.1],
[-0.9, 2.9],
[-0.5, 2.6],
[2.3, -1.1],
[2.7, -1.5]
])

y_train = torch.tensor([0, 0, 0, 1, 1])
X_test = torch.tensor([
[-0.8, 2.8],
[2.6, -1.6],
])
y_test = torch.tensor([0, 1])


from torch.utils.data import Dataset, DataLoader

class ToyDataset(Dataset):
    def __init__(self, x, y): 
        self.features = x 
        self.labels = y

    def __getitem__(self, index): 
        one_x = self.features[index]
        one_y = self.labels[index]
        return one_x, one_y

    def __len__(self):
        return self.labels.shape[0]

    
train_ds = ToyDataset(X_train, y_train)
test_ds = ToyDataset(X_test, y_test)

torch.manual_seed(123)
train_loader = DataLoader(
    dataset=train_ds,
    batch_size=2,
    shuffle=True,
    num_workers=0
)

test_loader = DataLoader(
    dataset=test_ds,
    batch_size=2,
    shuffle=False,
    num_workers=0
)


for idx, (x, y) in enumerate(train_loader):
    print(f"Batch {idx+1}: ", x, y)




