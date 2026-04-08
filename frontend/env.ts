import * as z from 'zod'

const coercedBoolean = z.string()
  .refine(s => s === 'true' || s === 'false' || s === '0' || s === '1')
  .transform(s => s === 'true' || s === '1')
  .optional()
const coercedNumber = z.coerce.number().int().positive().optional()
const optionalString = z.string().optional()

const envSchema = z.object({
  VITE_ALLOW_EMBED: coercedBoolean.default('false'),
  VITE_ALLOW_UNSAFE_DATA_SCHEME: coercedBoolean.default('false'),
  VITE_AMPLITUDE_API_KEY: optionalString,
  VITE_API_PREFIX: optionalString,
  VITE_BASE_PATH: z.string().default(''),
  VITE_BATCH_CONCURRENCY: coercedNumber.default('5'),
  VITE_COOKIE_DOMAIN: optionalString,
  VITE_CSP_WHITELIST: optionalString,
  VITE_DEPLOY_ENV: z.enum(['DEVELOPMENT', 'PRODUCTION', 'TESTING']).optional(),
  VITE_DISABLE_UPLOAD_IMAGE_AS_ICON: coercedBoolean.default('false'),
  VITE_EDITION: z.enum(['SELF_HOSTED', 'CLOUD']).default('SELF_HOSTED'),
  VITE_ENABLE_SINGLE_DOLLAR_LATEX: coercedBoolean.default('false'),
  VITE_ENABLE_WEBSITE_FIRECRAWL: coercedBoolean.default('true'),
  VITE_ENABLE_WEBSITE_JINAREADER: coercedBoolean.default('true'),
  VITE_ENABLE_WEBSITE_WATERCRAWL: coercedBoolean.default('false'),
  VITE_GITHUB_ACCESS_TOKEN: optionalString,
  VITE_INDEXING_MAX_SEGMENTATION_TOKENS_LENGTH: coercedNumber.default('4000'),
  VITE_IS_MARKETPLACE: coercedBoolean.default('false'),
  VITE_LOOP_NODE_MAX_COUNT: coercedNumber.default('100'),
  VITE_MAINTENANCE_NOTICE: optionalString,
  VITE_MARKETPLACE_API_PREFIX: optionalString,
  VITE_MARKETPLACE_URL_PREFIX: optionalString,
  VITE_MAX_ITERATIONS_NUM: coercedNumber.default('99'),
  VITE_MAX_PARALLEL_LIMIT: coercedNumber.default('10'),
  VITE_MAX_TOOLS_NUM: coercedNumber.default('10'),
  VITE_MAX_TREE_DEPTH: coercedNumber.default('50'),
  VITE_PUBLIC_API_PREFIX: optionalString,
  VITE_SENTRY_DSN: optionalString,
  VITE_SITE_ABOUT: optionalString,
  VITE_SUPPORT_EMAIL_ADDRESS: optionalString,
  VITE_SUPPORT_MAIL_LOGIN: coercedBoolean.default('false'),
  VITE_TEXT_GENERATION_TIMEOUT_MS: coercedNumber.default('60000'),
  VITE_TOP_K_MAX_VALUE: coercedNumber.default('10'),
  VITE_UPLOAD_IMAGE_AS_ICON: coercedBoolean.default('false'),
  VITE_WEB_PREFIX: optionalString,
  VITE_ZENDESK_FIELD_ID_EMAIL: optionalString,
  VITE_ZENDESK_FIELD_ID_ENVIRONMENT: optionalString,
  VITE_ZENDESK_FIELD_ID_PLAN: optionalString,
  VITE_ZENDESK_FIELD_ID_VERSION: optionalString,
  VITE_ZENDESK_FIELD_ID_WORKSPACE_ID: optionalString,
  VITE_ZENDESK_WIDGET_KEY: optionalString,
})

const parsed = envSchema.safeParse(import.meta.env)

if (!parsed.success) {
  console.error('Invalid environment variables:', parsed.error.flatten().fieldErrors)
}

export const env = parsed.success ? parsed.data : envSchema.parse({})
