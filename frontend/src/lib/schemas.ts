import { z } from "zod";

const ShortenUrlRequestSchema = z.object({
    originalUrl: z
        .string()
        .url("Please enter a valid URL"),
    ttl: z
        .number()
        .min(1, "TTL must be at least 1 minute")
        .max(60 * 24 * 365, "TTL cannot exceed 1 year"),
});

export type ShortenUrlRequest = z.infer<typeof ShortenUrlRequestSchema>;

const ShortenUrlResponseSchema = z.object({
    shortUrl: z.string().url(),
    expiresAt: z.string().refine(
        (date) => !isNaN(Date.parse(date)),
        {
            message: "Invalid date format",
        }
    ),
});

export type ShortenUrlResponse = z.infer<typeof ShortenUrlResponseSchema>;