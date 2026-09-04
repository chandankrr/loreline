import {
  Body,
  Container,
  Head,
  Html,
  Link,
  Preview,
  Section,
  Tailwind,
  Text,
} from "react-email";

import LorelineFonts from "../font.js";
import lorelineTailwindConfig from "../theme.js";

type EmailVerificationEmailProps = {
  baseUrl: string;
  verificationCode: string;
  expiresInMinutes?: number;
};

export default function EmailVerificationEmail({
  baseUrl,
  verificationCode,
  expiresInMinutes = 15,
}: EmailVerificationEmailProps) {
  const brand = "Loreline";

  return (
    <Tailwind config={lorelineTailwindConfig}>
      <Html>
        <Head>
          <LorelineFonts />
        </Head>
        <Preview>
          Verify your email to start reading with {brand} — code inside.
        </Preview>
        <Body className="m-0 p-0 font-sans">
          <Container className="mx-auto max-w-150 px-4 py-12">
            <Section className="overflow-hidden rounded-lg border border-border bg-surface shadow-card">
              {/* Hero */}
              <Section className="mobile:px-0 px-10 mobile:py-0 py-14">
                <Text className="font-40 font-display text-ink italic">
                  Confirm your email
                </Text>
                <Text className="m-0 mt-3 max-w-105 font-16 font-sans text-ink-soft">
                  One last step before you can start reading — confirm this is
                  your email address so Loreline knows where to keep your
                  library.
                </Text>
              </Section>

              {/* Code fallback */}
              <Section className="border-border border-t bg-surface-muted mobile:px-0! px-10 mobile:py-0 py-14">
                <Text className="m-0 max-w-95 font-22 font-display text-ink italic">
                  Your verification code
                </Text>
                <Text className="m-0 mt-2 max-w-105 font-14 font-sans text-ink-soft">
                  Enter the code below to confirm your email.
                </Text>

                <Section className="mt-6 rounded-lg border border-brand bg-brand-soft px-6 py-5 text-center">
                  <Text className="m-0 font-30 font-sans font-semibold text-brand-ink tracking-[8px]">
                    {verificationCode}
                  </Text>
                </Section>

                <Text className="m-0 mt-4 font-13 font-sans text-ink-faint">
                  This code expires in {expiresInMinutes} minutes.
                </Text>
              </Section>

              {/* Security footer */}
              <Section className="border-border border-t mobile:px-0 px-10 mobile:py-0 py-12">
                <Text className="m-0 max-w-105 font-13 font-sans text-ink-soft">
                  If you didn&apos;t create a {brand} account, you can safely
                  ignore this email — no account will be created without
                  verification.
                </Text>

                <Text className="mt-8 font-12 font-sans text-ink-faint">
                  123 Market Street, Floor 1
                  <br />
                  Tech City, IN, 941026
                </Text>

                <Text className="m-0 mt-4 max-w-[320px] font-12 font-sans text-ink-faint">
                  Need help?{" "}
                  <Link
                    href={`${baseUrl}/support`}
                    className="text-ink-faint underline"
                  >
                    Contact support
                  </Link>
                  .
                </Text>
              </Section>
            </Section>
          </Container>
        </Body>
      </Html>
    </Tailwind>
  );
}

EmailVerificationEmail.PreviewProps = {
  baseUrl: "https://example.com",
  verificationCode: "482913",
  expiresInMinutes: 15,
} satisfies EmailVerificationEmailProps;
