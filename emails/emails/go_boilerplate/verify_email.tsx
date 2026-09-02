// Go Boilerplate emails — Protocol activation (Figma node 2:7186).
// https://www.figma.com/design/vkeMjKNcfWdUvmfHdUfL8H/SaaS-Email-Templates--Community-?node-id=2-7186

import {
  Body,
  Button,
  Container,
  Head,
  Html,
  Img,
  Preview,
  Section,
  Tailwind,
  Text,
} from "react-email";
import { goBoilerplateAsset } from "./assets";
import { GoBoilerplateFonts } from "./go_boilerplate_fonts";
import { GoBoilerplateFooter } from "./go_boilerplate_footer";
import { GoBoilerplateHeader } from "./go_boilerplate_header";
import { GoBoilerplateOtpBlock } from "./go_boilerplate_otp_block";
import { goBoilerplateTailwindConfig } from "./theme";

export type VerifyEmailProps = {
  companyName?: string;
  name?: string;
  code: string;
  expiryMinutes?: number;
  verifyUrl?: string;
};

export function VerifyEmail({
  companyName = "Go Boilerplate",
  name = "there",
  code,
  expiryMinutes = 10,
  verifyUrl,
}: VerifyEmailProps) {
  return (
    <Tailwind config={goBoilerplateTailwindConfig}>
      <Html>
        <Head>
          <GoBoilerplateFonts />
        </Head>

        <Body className="bg-bg font-15 text-fg m-0 p-0 font-sans">
          <Preview>Confirm your {companyName} email address</Preview>
          <Container className="bg-bg mx-auto max-w-[640px]">
            <GoBoilerplateHeader companyName={companyName} />

            <Section align="left" className="px-4">
              <Img
                src={goBoilerplateAsset("hero_image.png")}
                alt=""
                width={608}
                height={240}
                className="block h-[240px] w-full max-w-[608px] border-none object-cover"
              />
            </Section>

            <Section align="left" className="mobile:px-4 mobile:py-10 px-10 py-14">
              <Section align="left" className="mobile:mb-8 mb-10">
                <Text className="font-56 mobile:font-40 text-fg m-0 font-condensed uppercase">
                  almost there
                </Text>
                <Text className="font-15 text-fg-2 m-0 mt-6 font-sans">
                  Hi {name}, thank you for signing up for {companyName}.
                </Text>
                <Text className="font-15 text-fg-2 m-0 mt-4 font-sans">
                  Enter this verification code to confirm your email, or use the
                  button below.
                </Text>
              </Section>

              <GoBoilerplateOtpBlock code={code} expiryMinutes={expiryMinutes} />

              {verifyUrl ? (
                <Button
                  href={verifyUrl}
                  className="bg-surface text-ink font-15-medium mt-10 box-border rounded-[1px] px-6 py-3 no-underline"
                >
                  Confirm Email
                </Button>
              ) : null}
            </Section>

            <Section align="left" className="mobile:px-4 px-10 pb-12">
              <Text className="font-11 text-fg-3 m-0 font-sans">
                If you didn&apos;t create an account, you can safely ignore this
                email.
              </Text>
            </Section>

            <GoBoilerplateFooter companyName={companyName} />
          </Container>
        </Body>
      </Html>
    </Tailwind>
  );
}

VerifyEmail.PreviewProps = {
  companyName: "Go Boilerplate",
  name: "Ada",
  code: "739204",
  expiryMinutes: 10,
  verifyUrl: "http://localhost:8080/api/v1/auth/verify-email?token=preview",
} satisfies VerifyEmailProps;

export default VerifyEmail;
