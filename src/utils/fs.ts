import fs from 'fs';
import path from 'path';
import fg from 'fast-glob';

export const createDirectory = async (directoryPath: string): Promise<void> => {
	try {
		await fs.promises.stat(directoryPath);
	} catch (err) {
		await fs.promises.mkdir(directoryPath);
	}
};

export const getDirectories = async (
	directoryPath: string,
	options?: {
		fileExtension?: string;
		recursive?: boolean;
	},
): Promise<string[]> => {
	const fileExtension = options?.fileExtension;
	const recursive = options?.recursive;

	const files = await getFiles(directoryPath, {
		fileExtension,
		recursive,
	});

	return [...new Set(files.map((directory) => path.parse(directory).dir))];
};

export const getFiles = async (
	directoryPath: string,
	options?: {
		fileExtension?: string;
		recursive?: boolean;
	},
): Promise<string[]> => {
	const fileExtension = options?.fileExtension;
	const recursive = options?.recursive;

	let pattern = fg.escapePath(directoryPath);

	if (recursive) {
		pattern = `${pattern}/**/*`;
	} else {
		pattern = `${pattern}/*`;
	}

	if (fileExtension) {
		pattern = `${pattern}.${fileExtension}`;
	}

	return await fg(pattern, {
		extglob: true,
	});
};
