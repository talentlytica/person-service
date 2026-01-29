/**
 * Extract passed and failed test cases from test-results.json
 * Writes: Passed_testcase.md and Failed_testcase.md in Test/
 */

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const resultsPath = path.join(__dirname, 'test-results.json');
const statusDir = path.join(__dirname, 'Test status');
const passedPath = path.join(statusDir, 'Passed_testcase.md');
const failedPath = path.join(statusDir, 'Failed_testcase.md');

const data = JSON.parse(fs.readFileSync(resultsPath, 'utf8'));

const passed = [];
const failed = [];

for (const suite of data.testResults || []) {
  const file = path.relative(path.join(__dirname), suite.testFilePath).replace(/\\/g, '/');
  for (const test of suite.testResults || []) {
    const entry = {
      file,
      suite: (test.ancestorTitles || []).join(' > '),
      title: test.title,
      fullName: test.fullName,
      duration: test.duration,
    };
    if (test.status === 'passed') {
      passed.push(entry);
    } else if (test.status === 'failed') {
      entry.failureMessages = (test.failureMessages || []).join('\n').trim();
      failed.push(entry);
    }
  }
}

function formatList(list, title, emoji) {
  let out = `# ${title}\n\n`;
  out += `**Total: ${list.length}**\n\n`;
  out += `| # | Test Name | File | Duration (ms) |\n`;
  out += `|---|-----------|------|---------------|\n`;
  list.forEach((t, i) => {
    const name = (t.title || t.fullName || '').replace(/\|/g, '\\|').replace(/\n/g, ' ');
    out += `| ${i + 1} | ${name} | \`${t.file}\` | ${t.duration != null ? t.duration : '-'} |\n`;
  });
  return out;
}

function formatFailedList(list) {
  let out = `# Failed Test Cases\n\n`;
  out += `**Total: ${list.length}**\n\n`;
  list.forEach((t, i) => {
    out += `## ${i + 1}. ${(t.title || t.fullName || '').replace(/\n/g, ' ')}\n\n`;
    out += `- **File:** \`${t.file}\`\n`;
    out += `- **Suite:** ${t.suite || '-'}\n`;
    if (t.duration != null) out += `- **Duration:** ${t.duration} ms\n`;
    if (t.failureMessages) {
      out += `\n**Failure message:**\n\`\`\`\n${t.failureMessages}\n\`\`\`\n\n`;
    }
  });
  return out;
}

if (!fs.existsSync(statusDir)) {
  fs.mkdirSync(statusDir, { recursive: true });
}
fs.writeFileSync(passedPath, formatList(passed, 'Passed Test Cases', '✅'), 'utf8');
fs.writeFileSync(failedPath, formatFailedList(failed), 'utf8');

console.log(`✅ Written ${passed.length} passed tests to Test status/Passed_testcase.md`);
console.log(`✅ Written ${failed.length} failed tests to Test status/Failed_testcase.md`);
