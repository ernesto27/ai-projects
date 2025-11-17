import re

with open("The_Verdict.txt", "r", encoding="utf-8") as file:
    content = file.read()

print("Total number of characters ", len(content))
# print(content[:99])

preprocessed = re.split(r'([,.:;?_!"()\']|--|\\s)', content)
preprocessed = [item.strip() for item in preprocessed if item.strip()]
print(len(preprocessed))
