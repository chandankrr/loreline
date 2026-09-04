import {
  Body,
  Button,
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

type PasswordResetEmailProps = {
  baseUrl: string;
  resetUrl: string;
  expiresInMinutes?: number;
};

export default function PasswordResetEmail({
  baseUrl,
  resetUrl,
  expiresInMinutes = 30,
}: PasswordResetEmailProps) {
  const brand = "Loreline";

  return (
    <Tailwind config={lorelineTailwindConfig}>
      <Html>
        <Head>
          <LorelineFonts />
        </Head>
        <Preview>Reset your {brand} password — this link expires soon.</Preview>
        <Body className="m-0 p-0 font-sans">
          <Container className="mx-auto max-w-150 px-4 py-12">
            <Section className="overflow-hidden rounded-lg border border-border bg-surface shadow-card">
              {/* Hero */}
              <Section className="mobile:px-0 px-10 mobile:py-0 py-14">
                <Text className="font-40 font-display text-ink italic">
                  Reset your password
                </Text>
                <Text className="m-0 mt-3 max-w-105 font-16 font-sans text-ink-soft">
                  We received a request to reset the password for your {brand}{" "}
                  account. Click below to choose a new one.
                </Text>

                <Section className="mt-8">
                  <Button
                    href={resetUrl}
                    className="inline-block rounded-full bg-brand px-6 py-2.5 text-center font-16 font-sans font-semibold text-brand-ink"
                  >
                    Reset password
                  </Button>
                </Section>

                <Text className="m-0 mt-4 font-13 font-sans text-ink-faint">
                  This link expires in {expiresInMinutes} minutes.
                </Text>
              </Section>

              {/* Security footer */}
              <Section className="border-border border-t mobile:px-0 px-10 mobile:py-0 py-12">
                <Text className="m-0 max-w-105 font-13 font-sans text-ink-soft">
                  If you didn&apos;t request this, you can safely ignore this
                  email — your password won&apos;t change unless you open the
                  link above and set a new one.
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

PasswordResetEmail.PreviewProps = {
  baseUrl: "https://example.com",
  resetUrl: "https://example.com/reset?token=abc123",
  expiresInMinutes: 30,
} satisfies PasswordResetEmailProps;
