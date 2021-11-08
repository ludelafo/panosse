import { executeCommand } from './exec';

export const getEncoderVersion = async () => {
	const stdout = await executeCommand(`flac --version`);

	const encoderVersionRegex = new RegExp('(flac) (d+[.]d+[.]d+)');

	return encoderVersionRegex.exec(stdout)?.at(1);
};

export const compressFile = async (flacFile: string, encoderSettings: string) =>
	await executeCommand(`${encoderSettings} "${flacFile}"`);
