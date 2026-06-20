import {
	createCipheriv,
	createDecipheriv,
	randomBytes,
} from "crypto";
import { BUANG_GHA_SECRET_KEY } from "../env.js";

const ALGO = "aes-256-gcm";
const IV_LEN = 12;
const TAG_LEN = 16;

const key = () => Buffer.from(BUANG_GHA_SECRET_KEY.value(), "hex");

// Returns: base64(iv || ciphertext || authTag)
export const encrypt = (plaintext: string): string => {
	const iv = randomBytes(IV_LEN);
	const cipher = createCipheriv(ALGO, key(), iv);
	const ct = Buffer.concat([cipher.update(plaintext, "utf8"), cipher.final()]);
	return Buffer.concat([iv, ct, cipher.getAuthTag()]).toString("base64");
};

export const decrypt = (encoded: string): string => {
	const buf = Buffer.from(encoded, "base64");
	const iv = buf.subarray(0, IV_LEN);
	const tag = buf.subarray(buf.length - TAG_LEN);
	const ct = buf.subarray(IV_LEN, buf.length - TAG_LEN);
	const decipher = createDecipheriv(ALGO, key(), iv);
	decipher.setAuthTag(tag);
	return Buffer.concat([decipher.update(ct), decipher.final()]).toString(
		"utf8",
	);
};
