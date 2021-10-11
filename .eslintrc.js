module.exports = {
	env: {
		node: true,
		es6: true,
	},
	ignorePatterns: ['!.eslintrc.js', '!.prettierrc.js'],
	extends: [
		'eslint:recommended',
		'plugin:@typescript-eslint/eslint-recommended',
		'plugin:@typescript-eslint/recommended',
		'plugin:prettier/recommended',
		'prettier',
	],
	plugins: ['prettier'],
	rules: {
		'@typescript-eslint/no-shadow': ['error'],
		'@typescript-eslint/no-unused-vars': ['error'],
		indent: ['error', 'tab'],
		'no-tabs': 0,
	},
};
