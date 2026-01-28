See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.

# The Ultimate Baseline Prompt for a Prompt Library

## Introduction to Prompt Engineering

Prompt engineering is the art and science of designing and refining inputs (prompts) to guide a generative AI model toward a desired output. It's a critical skill for getting the most out of large language models (LLMs).

This document, built on a foundation of expert prompt engineering principles and heavily inspired by the insights from NetworkChuck's "You SUCK at prompting AI (Here's the secret)" video, provides a comprehensive baseline for creating powerful and effective prompts. The core philosophy is that if the AI's response is poor, it's often a "personal skill issue"—the problem lies in the prompt itself. It's not just about asking a question; it's about programming the AI with words.

## Core Principles of Effective Prompts

These are the fundamental principles that separate a mediocre prompt from a great one.

### 1. Assign a Persona: Give the AI a Role

When you ask a generic question, you get a generic, soulless answer. By assigning a persona, you narrow the AI's focus and give it a specific perspective to draw from, dramatically improving the quality and style of the response.

*   **Less Effective:** "Write an apology email."
*   **More Effective:** "You are a senior site reliability engineer for CloudFlare, writing to both customers and engineers. Write an apology email for the recent outage."

### 2. Context is King: Be Clear and Specific

AI models will "hallucinate" or make things up to fill in gaps in their knowledge. The more context, background, and specific detail you provide, the less the AI has to guess. As NetworkChuck states, **"More context equals less hallucinations."**

*   **Less Effective:** "Tell me about cars."
*   **More Effective:** "I am a high school teacher creating a lesson plan on the American Revolution. Provide a summary of the key events leading up to the Declaration of Independence, suitable for 10th-grade students."
*   **Key Takeaway:** Be detailed and specific. Don't assume the AI knows anything. Always provide all the context, every time.

### 3. Give the AI Permission to Fail (Handling Ambiguity)

AIs are designed to be helpful and will often invent an answer rather than admit they don't know. To combat this, you must explicitly give the AI permission to fail. This is the number one fix for hallucinations.

*   **The Principle:** If the AI determines it lacks the necessary context, it must not invent information. Instead, it must summarize the context it has, state what's missing, and ask for the specific clarification it needs.
*   **Instruction:** "If you cannot find the answer in the provided context, you must say, 'I don't know.' Do not lie or invent information to please me."

### 4. Define the Output Format

Telling the LLM exactly how you want the result to look is a superpower. Don't leave the structure of the response to chance. Clearly specify the desired format, such as length, tone, and structure (e.g., list, table, JSON).

*   **Example:** "Create a table comparing the pros and cons of solar and wind energy. Keep the response under 200 words. The tone should be radically transparent, with no corporate fluff. Use a bulleted list for the timeline."

### 5. Use Examples (Few-Shot Prompting)

Instead of just *describing* the output, *show* the AI what you want. Providing one or more examples (known as one-shot or few-shot prompting) gives the model a clear pattern to follow and is one of the most effective ways to get the best results.

*   **Example:** "Translate the following English phrases to French in a formal tone:
    English: 'How are you?' -> French: 'Comment allez-vous?'
    English: 'Thank you very much.' -> French: 'Merci beaucoup.'
    English: 'Can you help me?' -> French:"

### 6. Use Affirmative Language

Tell the AI what to do, rather than what *not* to do. Positive and direct instructions are generally more effective and less ambiguous for the model to follow.

*   **Less Effective:** "Don't write a long response."
*   **More Effective:** "Summarize the article in three concise sentences."

## Advanced Prompting Techniques

Once you've mastered the basics, you can use these techniques for more complex reasoning tasks.

### Chain of Thought (COT)

Instruct the AI to "think step by step" before it answers. This forces the model to lay out its reasoning process, which improves accuracy and allows you to see *how* it arrived at a conclusion. Many AI platforms now have a "thinking" or "extended thinking" mode that automates this.

### Tree of Thought (TOT)

For even more complex problems, this approach has the AI explore multiple reasoning paths at once, like branches of a tree. It enables the AI to generate a diversity of options, evaluate different paths, and self-correct to find the best solution.

## The Meta-Skill: Clarity of Thought

The single most important skill in prompt engineering is **clarity of thought**. All the techniques are ultimately about one thing: expressing yourself clearly. If your thinking is messy, your prompts will be messy, and your results will be poor. The AI can only be as clear as you are.

**The Golden Rule: Think first, prompt second.** Before you write a prompt, take the time to clarify your own goals. If you can't explain it clearly to yourself, you can't prompt it effectively.

## Baseline Prompt Template

This template incorporates all the principles discussed, providing a robust foundation for your prompt library.

