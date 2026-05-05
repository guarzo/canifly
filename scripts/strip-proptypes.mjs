#!/usr/bin/env node
// Strip `import PropTypes from 'prop-types'` and `<Identifier>.propTypes = { ... };`
// blocks from React source files. Idempotent: running twice produces no further diffs.
import { readFileSync, writeFileSync } from 'node:fs';
import { argv } from 'node:process';

const files = argv.slice(2);
if (files.length === 0) {
    console.error('Usage: node scripts/strip-proptypes.mjs <file>...');
    process.exit(1);
}

let totalImportLinesRemoved = 0;
let totalBlocksRemoved = 0;

for (const file of files) {
    const before = readFileSync(file, 'utf8');
    let after = before;

    // 1) Strip the import line. Allow single or double quotes; allow trailing semicolon
    //    or none; consume up to and including the trailing newline so we don't leave a
    //    blank line behind.
    const importPattern = /^import\s+PropTypes\s+from\s+['"]prop-types['"];?\s*\r?\n/m;
    if (importPattern.test(after)) {
        after = after.replace(importPattern, '');
        totalImportLinesRemoved += 1;
    }

    // 2) Strip the `<Identifier>.propTypes = { ... };` block.
    //    The block is multi-line. We match from `<Word>.propTypes = {` through the
    //    matching closing `};` and consume any trailing blank line.
    //    JS regex doesn't have native brace-balancing, but in this codebase every
    //    propTypes block ends with `};` on its own line. Match non-greedy to the first
    //    `};` at the start of a line.
    const blockPattern = /\n[A-Za-z_$][\w$]*\.propTypes\s*=\s*\{[\s\S]*?\n\};\s*\r?\n/g;
    const blockMatches = after.match(blockPattern) || [];
    if (blockMatches.length > 0) {
        after = after.replace(blockPattern, '\n');
        totalBlocksRemoved += blockMatches.length;
    }

    if (after !== before) {
        writeFileSync(file, after, 'utf8');
        console.log(`stripped ${file}`);
    }
}

console.log(`done — removed ${totalImportLinesRemoved} import lines, ${totalBlocksRemoved} propTypes blocks`);
