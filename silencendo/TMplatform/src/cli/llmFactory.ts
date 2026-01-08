import { LLMClient } from '../llm/types';
import { MockLLM } from '../llm/mockLLM';
import { DeepSeekLLM } from '../llm/deepseek';

export function createLLMClient(): LLMClient {
  const useDeepSeek = process.env.DEEPSEEK_API_KEY !== undefined;
  
  if (useDeepSeek) {
    console.log('🚀 Using DeepSeek LLM');
    try {
      return new DeepSeekLLM();
    } catch (error) {
      console.warn(`⚠️  Failed to initialize DeepSeek: ${error}`);
      console.log('🔄 Falling back to Mock LLM');
      return new MockLLM();
    }
  } else {
    console.log('🧪 Using Mock LLM (set DEEPSEEK_API_KEY to use DeepSeek)');
    return new MockLLM();
  }
}