import { Column, Img, Row, Section, Text } from "react-email";
import { goBoilerplateAsset } from "./assets";

type GoBoilerplateHeaderProps = {
  companyName: string;
};

export function GoBoilerplateHeader({ companyName }: GoBoilerplateHeaderProps) {
  return (
    <Section align="left" className="mobile:px-4 px-10 py-6">
      <Row>
        <Column className="w-[23px] align-middle">
          <Img
            src={goBoilerplateAsset("logo_mark.svg")}
            alt=""
            width={23}
            height={23}
            className="block border-none"
          />
        </Column>
        <Column className="align-middle pl-2">
          <Text className="font-17 text-fg m-0 font-sans">{companyName}</Text>
        </Column>
      </Row>
    </Section>
  );
}
