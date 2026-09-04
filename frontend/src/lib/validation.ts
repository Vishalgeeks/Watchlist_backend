import { z } from 'zod';

export const registerSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters').max(50, 'Name must be at most 50 characters')
    .refine(val => val.trim() !== '', 'Name cannot be only spaces')
    .refine(val => /[a-zA-Z0-9]/.test(val), 'Name must contain at least one letter or number'),
  email: z.string().email('Invalid email format'),
  password: z.string().min(8, 'Password must be at least 8 characters')
    .refine(val => /[A-Z]/.test(val), 'Password must contain an uppercase letter')
    .refine(val => /[a-z]/.test(val), 'Password must contain a lowercase letter')
    .refine(val => /[0-9]/.test(val), 'Password must contain a number')
    .refine(val => /[^A-Za-z0-9]/.test(val), 'Password must contain a special character'),
  confirmPassword: z.string(),
}).refine(data => data.password === data.confirmPassword, {
  message: 'Passwords do not match',
  path: ['confirmPassword'],
});

export type RegisterInput = z.infer<typeof registerSchema>;

export const loginSchema = z.object({
  email: z.string().email('Invalid email format'),
  password: z.string().min(1, 'Password is required'),
});

export type LoginInput = z.infer<typeof loginSchema>;

export const createWatchlistSchema = z.object({
  name: z.string().min(1, 'Name is required').max(100, 'Name must be at most 100 characters')
    .refine(val => val.trim() !== '', 'Name cannot be only spaces')
    .refine(val => /[a-zA-Z0-9]/.test(val), 'Name must contain at least one letter or number'),
});

export type CreateWatchlistInput = z.infer<typeof createWatchlistSchema>;

export const addStockSchema = z.object({
  stock_id: z.number().int().min(1, 'Stock ID is required'),
});

export type AddStockInput = z.infer<typeof addStockSchema>;

export const createOrderSchema = z.object({
  stock_id: z.number().int().min(1, 'Stock is required'),
  side: z.enum(['BUY', 'SELL']),
  order_type: z.enum(['MARKET', 'LIMIT']),
  quantity: z.number().int().min(1, 'Quantity must be at least 1'),
  price: z.number().positive('Price must be greater than 0').optional(),
}).refine(data => {
  if (data.order_type === 'LIMIT' && (data.price === undefined || data.price <= 0)) {
    return false;
  }
  return true;
}, {
  message: 'Price is required and must be greater than 0 for LIMIT orders',
  path: ['price'],
});

export type CreateOrderInput = z.infer<typeof createOrderSchema>;

export const depositSchema = z.object({
  amount: z.number().positive('Amount must be greater than 0'),
  description: z.string().optional(),
});

export type DepositInput = z.infer<typeof depositSchema>;

export const withdrawSchema = z.object({
  amount: z.number().positive('Amount must be greater than 0'),
  description: z.string().optional(),
});

export type WithdrawInput = z.infer<typeof withdrawSchema>;

export function formatValidationErrors(errors: z.ZodError): Record<string, string> {
  const fieldErrors: Record<string, string> = {};
  errors.errors.forEach(err => {
    const path = err.path.join('.');
    if (!fieldErrors[path]) {
      fieldErrors[path] = err.message;
    }
  });
  return fieldErrors;
}