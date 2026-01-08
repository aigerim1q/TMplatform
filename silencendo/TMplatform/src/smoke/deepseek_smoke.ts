import 'dotenv/config';
import { DeepSeekLLM } from '../llm/deepseek';
import { Chunk } from '../ingestion/types';
import { Context } from '../context/types';

async function smokeTest() {
  console.log('🧪 Starting DeepSeek smoke test...');
  
  try {
    // Initialize DeepSeek client
   const llm = new DeepSeekLLM();
    console.log('✅ DeepSeek runtime connected');
    
    // Test basic generation
    console.log('\n📝 Testing basic generation...');
    const testMessages = [
      { role: 'user' as const, content: 'Say "DeepSeek is working" in one sentence.' }
    ];
    
    const response = await llm.generate(testMessages);
    console.log('Response:', response);
    
    // Test answer functionality with mock data
    console.log('\n🔍 Testing answer functionality...');
    const mockChunks: Chunk[] = [
      {
        chunkId: 'test-1',
        sourceId: 'source-1',
        text: 'Artificial Intelligence and Machine Learning are transforming industries.',
        metadata: { position: 0, length: 100 }
      }
    ];
    
    const mockContext: Context = {
      project: 'Test Project',
      stage: 'Analysis',
      task: 'Testing',
      goals: ['Test connectivity'],
      constraints: ['API rate limits'],
      decisions: ['Use DeepSeek'],
      next_steps: ['Run full tests']
    };
    
    const answerResponse = await llm.answer('What is this document about?', mockChunks, mockContext);
    console.log('Answer response:', answerResponse.substring(0, 200) + '...');
    
    // Test edit functionality
    console.log('\n✏️  Testing edit functionality...');
    const editResponse = await llm.edit(
      'This is a long text that should be shortened.', 
      'Make this shorter', 
      mockChunks, 
      mockContext
    );
    console.log('Edit response:', editResponse.substring(0, 200) + '...');
    
    console.log('\n🎉 All smoke tests passed! DeepSeek integration is working correctly.');
    
  } catch (error) {
    console.error('❌ Smoke test failed:', error);
    process.exit(1);
  }
}

if (require.main === module) {
  smokeTest();
}