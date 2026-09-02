import { Column, Img, Link, Row, Section, Text } from "react-email";
import { goBoilerplateAsset } from "./assets";

type GoBoilerplateFooterProps = {
  companyName: string;
  tagline?: string;
};

const socialLinks = [
  { href: "https://x.com/", icon: "social_x.svg", alt: "X" },
  { href: "https://linkedin.com/", icon: "social_linkedin.svg", alt: "LinkedIn" },
  { href: "https://youtube.com/", icon: "social_youtube.svg", alt: "YouTube" },
  {
    href: "https://github.com/",
    icon: "social_github.svg",
    alt: "GitHub",
  },
] as const;

export function GoBoilerplateFooter({
  companyName,
  tagline = "Go Boilerplate is a production-ready Go HTTP API starter with auth, email, and OpenAPI.",
}: GoBoilerplateFooterProps) {
  return (
    <Section align="left" className="border-stroke mobile:px-4 mobile:py-12 border-t px-10 py-16">
      <Text className="font-13 text-fg-3 m-0 max-w-[320px] font-sans">
        {tagline}
      </Text>

      <Row align="left">
        <Column className="w-full align-top">
          <Section align="left" className="mt-8 w-[152px]">
            <Row align="left">
              {socialLinks.map((social) => (
                <Column key={social.alt} className="w-[20px] pr-8">
                  <Link href={social.href} className="inline-block">
                    <Img
                      src={goBoilerplateAsset(social.icon)}
                      alt={social.alt}
                      width={20}
                      height={20}
                      className="block border-none"
                    />
                  </Link>
                </Column>
              ))}
            </Row>
          </Section>
        </Column>
      </Row>

      <Row align="left">
        <Column className="w-full pt-5 align-top">
          <Text className="font-11 text-fg-3 m-0 max-w-[169px] font-sans">
            Transactional email from {companyName}.
          </Text>
        </Column>
      </Row>
    </Section>
  );
}
