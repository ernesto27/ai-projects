from dotenv import load_dotenv

load_dotenv()

from anthropic import Anthropic

client = Anthropic()
model = "claude-sonnet-4-0"


def chat(messages, system=None):
    params = {
        "model": model,
        "max_tokens": 1000,
        "messages": messages,
    }

    if system:
        params["system"] = system


    message = client.messages.create(**params)
    return message.content[0].text


def add_user_messages(messages, text):
    user_message = {"role": "user", "content": text}
    messages.append(user_message)

def add_assistant_messages(messages, text):
    assistant_message = {"role": "assistant", "content": text}
    messages.append(assistant_message)

system = """
You are a patient math tutor.
Do not directly answer the question.
Guide them to a solution step by step.
"""

messages = [] 
add_user_messages(messages, "How do i solve 5x+3=2 for x?")
anwser = chat(messages)
print(anwser)

# add_assistant_messages(messages, anwser)

# add_user_messages(messages, "go more deep")
# anwser = chat(messages)
# print(anwser)


