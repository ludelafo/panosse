import path from 'path';
import { Worker } from 'worker_threads';
import { exec } from 'child_process';

export const splitFiles = (
	tasks: string[],
	threadCount: number,
): string[][] => {
	const tasksPerThread = Math.ceil(tasks.length / threadCount);

	return [...Array(threadCount)].map((_, thread) =>
		tasks.slice(
			thread * tasksPerThread,
			thread * tasksPerThread + tasksPerThread,
		),
	);
};

export const runThreads = async <T>(
	filesPerThread: string[][],
	options: T,
	tsWorkerPath: string,
	logger: Console,
): Promise<void> =>
	new Promise((resolve) => {
		const threads = new Set();

		filesPerThread.forEach((files, thread) => {
			const worker = new Worker(path.join(__dirname, './worker.js'), {
				workerData: {
					path: tsWorkerPath,
					threadNumber: thread,
					files,
					options,
				},
			});

			threads.add(worker);

			worker.on('message', (message) => logger.log(message));

			worker.on('exit', () => {
				threads.delete(worker);

				if (threads.size === 0) {
					resolve();
				}
			});
		});
	});

export const executeCommand = async (command: string): Promise<string> =>
	new Promise((resolve, reject) => {
		exec(command, (error, stdout, stderr) => {
			if (error) {
				reject(error);
			}

			if (stderr) {
				reject(stderr);
			}

			resolve(stdout);
		});
	});
