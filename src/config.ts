export const PORT = parseInt(process.env.PORT || '3000');
export const POSTGRES_HOST = process.env.POSTGRES_HOST || 'localhost';
export const POSTGRES_USER = process.env.POSTGRES_USER || 'postgres';
export const POSTGRES_PASSWORD = process.env.POSTGRES_PASSWORD || 'postgres';
export const POSTGRES_DATABASE = process.env.POSTGRES_DATABASE || POSTGRES_USER;
