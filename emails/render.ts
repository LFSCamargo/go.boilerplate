import * as React from "react";
import { render } from "react-email";
import PasswordResetEmail from "./emails/go_boilerplate/password_reset";
import VerifyEmail from "./emails/go_boilerplate/verify_email";

const templates: Record<string, React.ComponentType<any>> = {
  verify_email: VerifyEmail,
  password_reset: PasswordResetEmail,
};

async function main() {
  const templateName = process.argv[2];
  const propsJson = process.argv[3] ?? "{}";

  if (!templateName) {
    console.error("Usage: npm run render -- <template> '<json-props>'");
    process.exit(1);
  }

  const Component = templates[templateName];
  if (!Component) {
    console.error(`Unknown template: ${templateName}`);
    console.error(`Available: ${Object.keys(templates).join(", ")}`);
    process.exit(1);
  }

  let props: Record<string, unknown> = {};
  try {
    props = JSON.parse(propsJson);
  } catch {
    console.error("Invalid JSON props");
    process.exit(1);
  }

  const html = await render(React.createElement(Component, props));
  process.stdout.write(html);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
