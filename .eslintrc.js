module.exports = {
	root: true,
	env: {
		jest: true,
		node: true,
	},
	parser: '@typescript-eslint/parser',
	parserOptions: {
		project: 'tsconfig.json',
		sourceType: 'module',
	},
	ignorePatterns: ['!.eslintrc.js', '!.prettierrc.js'],
	extends: [
    'oclif',
    'oclif-typescript',
		'eslint:recommended',
		'plugin:@typescript-eslint/eslint-recommended',
		'plugin:@typescript-eslint/recommended',
		'plugin:prettier/recommended',
	],
	plugins: ['@typescript-eslint/eslint-plugin', 'prettier'],
	rules: {
		'@typescript-eslint/no-shadow': ['error'],
		'@typescript-eslint/no-unused-vars': ['error'],
	},
};
