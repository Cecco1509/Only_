import type { UserInformation } from '$lib/User/user';
import type { RequestEvent } from '@sveltejs/kit';

process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';

export const authenticateUser = async (event: RequestEvent): Promise<UserInformation | null> => {
	const token: string | undefined = event.cookies.get('token');
	console.log('Token string:', token);
	if (!token) {
		console.log('No token found');
		return null;
	}

	console.log('Token:', token);

	try {
		//https://authservice:5000
		const response = await fetch('http://localhost/api/auth/verify', {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			}
		});

		if (response.status === 200) {
			console.log('User authenticated');
			try {
				const json: { Status: number; Meta: object; Data: { [key: string]: unknown } } =
					await response.json();
				return {
					username: json.Data['username'] as string,
					email: json.Data['email'] as string,
					id: json.Data['id'] as number,
					token: token
				};
			} catch (error) {
				console.error(error);
				return null;
			}
		} else {
			console.log('User not authenticated');
			return null;
		}
	} catch (error) {
		console.error(error);
		return null;
	}
};
