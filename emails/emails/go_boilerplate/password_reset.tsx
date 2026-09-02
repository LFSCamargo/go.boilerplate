// Go Boilerplate emails — Protocol password reset (Figma node 2:7395).
// https://www.figma.com/design/vkeMjKNcfWdUvmfHdUfL8H/SaaS-Email-Templates--Community-?node-id=2-7395

import {
  Body,
  Button,
  Container,
  Head,
  Html,
  Preview,
  Section,
  Tailwind,
  Text,
} from "react-email";
import { GoBoilerplateFonts } from "./go_boilerplate_fonts";
import { GoBoilerplateFooter } from "./go_boilerplate_footer";
import { GoBoilerplateHeader } from "./go_boilerplate_header";
import { GoBoilerplateOtpBlock } from "./go_boilerplate_otp_block";
import { goBoilerplateTailwindConfig } from "./theme";

export type PasswordResetEmailProps = {
  companyName?: string;
  name?: string;
  code: string;
  expiryMinutes?: number;
  resetUrl?: string;
};

export function PasswordResetEmail({
  companyName = "Go Boilerplate",
  name = "there",
  code,
  expiryMinutes = 10,
  resetUrl,
}: PasswordResetEmailProps) {
  return (
    <Tailwind config={goBoilerplateTailwindConfig}>
      <Html>
        <Head>
          <GoBoilerplateFonts />
        </Head>

        <Body className="bg-bg font-15 text-fg m-0 p-0 font-sans">
          <Preview>Reset your {companyName} password</Preview>
          <Container className="bg-bg mx-auto max-w-[640px]">
            <GoBoilerplateHeader companyName={companyName} />

            <Section align="left" className="mobile:px-4 mobile:py-10 px-10 py-14">
              <Section align="left" className="mobile:mb-8 mb-10">
                <Text className="font-56 mobile:font-40 text-fg m-0 font-condensed uppercase">
                  Password reset
                </Text>
                <Text className="font-15 text-fg-2 m-0 mt-6 font-sans">
                  Hi {name}, we received a request to reset the password for
                  your {companyName} account. Use the code below to choose a new
                  password.
                </Text>
              </Section>

              <GoBoilerplateOtpBlock
                code={code}
                expiryMinutes={expiryMinutes}
                hint="If you didn't request this, please ignore this email. Your password won't change until you complete the reset flow."
              />

              {resetUrl ? (
                <Button
                  href={resetUrl}
                  className="bg-surface text-ink font-15-medium mt-10 box-border rounded-[1px] px-6 py-3 no-underline"
                >
                  Reset Password
                </Button>
              ) : null}
            </Section>

            <GoBoilerplateFooter companyName={companyName} />
          </Container>
        </Body>
      </Html>
    </Tailwind>
  );
}

PasswordResetEmail.PreviewProps = {
  companyName: "Go Boilerplate",
  name: "Ada",
  code: "482913",
  expiryMinutes: 10,
  resetUrl: "http://localhost:8080/api/v1/auth/reset-password",
} satisfies PasswordResetEmailProps;

export default PasswordResetEmail;
