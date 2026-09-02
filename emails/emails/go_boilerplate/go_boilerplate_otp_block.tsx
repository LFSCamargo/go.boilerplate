import { Section, Text } from "react-email";

type GoBoilerplateOtpBlockProps = {
  code: string;
  expiryMinutes: number;
  hint?: string;
};

const otpBoxStyle = {
  width: "100%",
  maxWidth: "100%",
  backgroundColor: "#212121",
  border: "1px solid #2B2B2B",
  borderRadius: "1px",
  padding: "16px 24px",
  textAlign: "center" as const,
  boxSizing: "border-box" as const,
};

export function GoBoilerplateOtpBlock({
  code,
  expiryMinutes,
  hint,
}: GoBoilerplateOtpBlockProps) {
  return (
    <>
      <Section align="left" style={otpBoxStyle}>
        <Text
          style={{ margin: 0, textAlign: "center" }}
          className="font-56 text-fg font-condensed tracking-[8px] uppercase"
        >
          {code}
        </Text>
      </Section>
      <Section align="left" style={{ width: "100%" }}>
        <Text className="font-11 text-fg-3 m-0 mt-4 font-sans">
          This code expires in {expiryMinutes} minutes.
        </Text>
        {hint ? (
          <Text className="font-11 text-fg-3 m-0 mt-2 font-sans">{hint}</Text>
        ) : null}
      </Section>
    </>
  );
}
