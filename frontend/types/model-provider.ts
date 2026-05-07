/**
 * Model provider quota types - shared type definitions for API responses
 * These represent the provider identifiers that support paid/trial quotas
 */
export enum ModelProviderQuotaGetPaid {
  ANTHROPIC = 'odysseythink/anthropic/anthropic',
  OPENAI = 'odysseythink/openai/openai',
  // AZURE_OPENAI = 'odysseythink/azure_openai/azure_openai',
  GEMINI = 'odysseythink/gemini/google',
  X = 'odysseythink/x/x',
  DEEPSEEK = 'odysseythink/deepseek/deepseek',
  TONGYI = 'odysseythink/tongyi/tongyi',
}
