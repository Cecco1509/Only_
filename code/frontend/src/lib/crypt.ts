import { base64ToBytes } from 'byte-base64';

export const deriveAESCryptoKeyFromPassword = async (
	passWordCryptoKey: CryptoKey,
	salt: BufferSource
) => {
	const cryptoKey = await window.crypto.subtle.deriveKey(
		{
			name: 'PBKDF2',
			salt,
			iterations: 100000,
			hash: 'SHA-256'
		},
		passWordCryptoKey,
		{ name: 'AES-GCM', length: 128 },
		false,
		['encrypt', 'decrypt']
	);
	return cryptoKey;
};

const encode = (text: string) => {
	return new TextEncoder().encode(text);
};

const generatePBKDF2CryptoKey = async (password: string) => {
	const passWordAsBytes = encode(password);
	const passWordCryptoKey = await window.crypto.subtle.importKey(
		'raw',
		passWordAsBytes,
		'PBKDF2',
		false,
		['deriveKey']
	);
	return passWordCryptoKey;
};

export async function encrypt(buffer: ArrayBuffer, password: string): Promise<Uint8Array> {
	const passWordCryptoKey: CryptoKey = await generatePBKDF2CryptoKey(password);
	const salt: Uint8Array = crypto.getRandomValues(new Uint8Array(16));
	const cryptoKey: CryptoKey = await deriveAESCryptoKeyFromPassword(passWordCryptoKey, salt);
	const iv: Uint8Array = crypto.getRandomValues(new Uint8Array(16));
	const encryptedBuffer = await window.crypto.subtle.encrypt(
		{
			name: 'AES-GCM',
			iv
		},
		cryptoKey,
		buffer
	);
	const encryptedData: Uint8Array = new Uint8Array(
		encryptedBuffer.byteLength + iv.byteLength + salt.byteLength
	);
	encryptedData.set(salt, 0);
	encryptedData.set(iv, salt.byteLength);
	encryptedData.set(new Uint8Array(encryptedBuffer), salt.byteLength + iv.byteLength);

	return encryptedData;
}

export async function decrypt(encryptedData: string, password: string) {
	const encryptedBytes = base64ToBytes(encryptedData);
	const salt = encryptedBytes.slice(0, 16);
	const iv = encryptedBytes.slice(16, 32);
	const encryptedBuffer = encryptedBytes.slice(32);

	const passWordCryptoKey = await generatePBKDF2CryptoKey(password);
	const cryptoKey = await deriveAESCryptoKeyFromPassword(passWordCryptoKey, salt);
	const decryptedBuffer = await window.crypto.subtle.decrypt(
		{
			name: 'AES-GCM',
			iv
		},
		cryptoKey,
		encryptedBuffer
	);

	console.log('Decrypted Blob: ', new TextDecoder().decode(decryptedBuffer));

	return decryptedBuffer;
}
