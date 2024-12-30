import type { EncFile } from '$lib/EncFile/EncFile';
import type { UserInformation } from '$lib/User/user.js';
import type { LayoutServerLoad } from './$types';

type SharedWith = {
	shared_with_user_id: number;
	shared_with_username: string;
	shared_by_user_id: number;
	shared_by_username: string;
};

export const load: LayoutServerLoad = async ({
	locals
}): Promise<{
	user: UserInformation;
	ownedFiles: EncFile[];
	sharingUsers: SharedWith[];
}> => {
	console.log('Layout server load:', locals.user);

	try {
		// https://archiveservice:5000/
		const userFiles = fetch('http://localhost/api/archive/efile', {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${locals.user.token}`
			}
		});

		const usersSharingWithMe = fetch('http://localhost/api/archive/shared/users', {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${locals.user.token}`
			}
		});
		//const files = await userFiles.json();
		//console.log('Files:', files.Data);

		return {
			user: locals.user,
			ownedFiles: (await (await userFiles).json()).Data as EncFile[],
			sharingUsers: (await (await usersSharingWithMe).json()).Data as SharedWith[]
		};
	} catch (error) {
		console.error(error);
		return {
			user: locals.user,
			ownedFiles: [],
			sharingUsers: []
		};
	}

	// shared people list
};
