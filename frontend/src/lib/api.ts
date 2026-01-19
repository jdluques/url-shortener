import { env } from "./env";
import { ShortenUrlRequest, ShortenUrlResponse } from "./schemas";

export async function shortenUrl(data: ShortenUrlRequest): Promise<ShortenUrlResponse> {
    const response = await fetch(`${env.API_BASE_URL}/shorten`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Failed to shorten URL");
    }

    return response.json();
}