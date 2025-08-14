#!/usr/bin/env python3

from dotenv import load_dotenv
import os
from groq import Groq

# Load environment variables
load_dotenv()
GROQ_API_KEY = os.getenv('GROQ_API_KEY')

print(f"API Key loaded: {'Yes' if GROQ_API_KEY else 'No'}")
if GROQ_API_KEY:
    print(f"API Key starts with: {GROQ_API_KEY[:10]}...")

try:
    client = Groq(api_key=GROQ_API_KEY)
    print("Groq client created successfully")
    
    # Test a simple chat completion
    chat_completion = client.chat.completions.create(
        messages=[
            {
                "role": "user",
                "content": "Hello, can you respond with just 'Hello back'?"
            }
        ],
        model="llama3-8b-8192",
    )
    
    response = chat_completion.choices[0].message.content
    print(f"Groq response: {response}")
    print("Test successful!")
    
except Exception as e:
    print(f"Error testing Groq: {str(e)}")
    import traceback
    traceback.print_exc()
