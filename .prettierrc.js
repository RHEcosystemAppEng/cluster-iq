export default {
  // Use single quotes instead of double quotes
  // Example: const name = 'John' vs "John"
  singleQuote: true,

  // Maximum line length before wrapping
  // Longer lines will be broken into multiple lines
  printWidth: 120,

  // Add semicolons at the end of statements
  // Example: const name = 'John';
  semi: true,

  // Add trailing commas in objects/arrays where valid in ES5
  // Example: { name: 'John', age: 30, }
  trailingComma: 'es5',

  // Number of spaces for each indentation level
  // Example:
  // {
  //   name: 'John'
  // }
  tabWidth: 2,

  // Add spaces between brackets in object literals
  // Example: { name: 'John' } vs {name: 'John'}
  bracketSpacing: true,

  // Put the closing bracket of JSX elements on a new line
  // Example:
  // <button
  //   className="btn"
  // >
  bracketSameLine: false,

  // Omit parentheses around single arrow function parameters
  // Example: x => x * 2 vs (x) => x * 2
  arrowParens: 'avoid',

  // Maintain existing line endings in your project
  // (Helps avoid issues between Windows/Unix systems)
  endOfLine: 'auto',

  // Special formatting rules for specific file types
  overrides: [
    {
      files: '*.md',
      options: {
        printWidth: 70,
        useTabs: false,
        trailingComma: 'none',
        arrowParens: 'avoid',
        proseWrap: 'never',
      },
    },
    {
      files: '*.{json,babelrc,eslintrc,remarkrc,prettierrc}',
      options: {
        useTabs: false,
      },
    },
  ],
};
