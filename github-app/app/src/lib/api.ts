const request = async <T>(
	path: string,
	options?: RequestInit,
): Promise<T> => {
	const res = await fetch(`/api${path}`, {
		headers: { "Content-Type": "application/json" },
		credentials: "same-origin",
		...options,
	});
	if (!res.ok) {
		const body = await res.json().catch(() => ({}));
		throw Object.assign(new Error("API error"), { status: res.status, body });
	}
	return res.json() as Promise<T>;
};

export const get = <T>(path: string) => request<T>(path);

export const post = <T>(path: string, data: unknown) =>
	request<T>(path, { method: "POST", body: JSON.stringify(data) });

export const patch = <T>(path: string, data: unknown) =>
	request<T>(path, { method: "PATCH", body: JSON.stringify(data) });
