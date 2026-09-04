import {
  Body,
  Button,
  Column,
  Container,
  Head,
  Html,
  Link,
  Preview,
  Row,
  Section,
  Tailwind,
  Text,
} from "react-email";

import LorelineFonts from "../font.js";
import lorelineTailwindConfig from "../theme.js";

type WelcomeEmailProps = {
  userFirstName: string;
  baseUrl: string;
};

const steps = [
  {
    number: "01",
    title: "Point",
    description:
      "Hover a word, a line, a diagram, or a whole paragraph. Loreline knows exactly what you mean — no copying, no describing where you are.",
  },
  {
    number: "02",
    title: "Ask",
    description:
      "Speak naturally. Interrupt, follow up, change direction, or just ask for the simpler version.",
  },
  {
    number: "03",
    title: "See",
    description:
      "Get an answer in voice while the sideboard draws the idea, the scene, or the relationship beside you.",
  },
];

export default function WelcomeEmail({
  userFirstName,
  baseUrl,
}: WelcomeEmailProps) {
  const brand = "Loreline";
  const welcomeTitle = `Welcome to ${brand}`;

  return (
    <Tailwind config={lorelineTailwindConfig}>
      <Html>
        <Head>
          <LorelineFonts />
        </Head>
        <Preview>{welcomeTitle}</Preview>
        <Body className="m-0 p-0 font-sans">
          <Container className="mx-auto max-w-150 px-4 py-12">
            <Section className="overflow-hidden rounded-lg border border-border bg-surface shadow-card">
              {/* Hero */}
              <Section className="mobile:px-0 px-10 mobile:py-0 py-14">
                <Text className="font-40 font-display text-ink italic">
                  {welcomeTitle}, {userFirstName}
                </Text>
                <Text className="m-0 mt-3 max-w-105 font-16 font-sans text-ink-soft">
                  You&apos;re set up. Bring any book, point at a line, and ask
                  what&apos;s on your mind — Loreline stays with the page and
                  answers by voice.
                </Text>

                <Section className="mt-8">
                  <Button
                    href={baseUrl}
                    className="inline-block rounded-full bg-brand px-6 py-2.5 text-center font-16 font-sans font-semibold text-brand-ink"
                  >
                    Open your first book
                  </Button>
                </Section>
              </Section>

              {/* How it works */}
              <Section className="border-border border-t bg-surface-muted mobile:px-0! px-10 mobile:py-0 py-14">
                <Text className="m-0 max-w-95 font-22 font-display text-ink italic">
                  The page is the prompt
                </Text>
                <Text className="m-0 mt-2 max-w-105 font-14 font-sans text-ink-soft">
                  No copying passages, no describing where you are. Three steps,
                  every time.
                </Text>

                <Section className="mt-8">
                  {steps.map((step, idx) => (
                    <Section
                      key={step.number}
                      className={
                        idx < steps.length - 1
                          ? "border-border border-b py-5"
                          : "py-5"
                      }
                    >
                      <Row>
                        <Column className="w-13 align-top">
                          <Text className="m-0 font-22 font-display text-ink-faint italic">
                            {step.number}
                          </Text>
                        </Column>
                        <Column className="align-top">
                          <Text className="m-0 font-16 font-sans font-semibold text-ink">
                            {step.title}
                          </Text>
                          <Text className="m-0 mt-1 max-w-95 font-14 font-sans text-ink-soft">
                            {step.description}
                          </Text>
                        </Column>
                      </Row>
                    </Section>
                  ))}
                </Section>
              </Section>

              {/* Epigraph */}
              <Section className="border-border border-t mobile:px-0 px-10 mobile:py-0 py-12">
                <Text className="m-0 font-18 font-display text-ink italic leading-relaxed">
                  &ldquo;Why does the author call attention a current?&rdquo;
                  <br />
                  &ldquo;Because a current already has direction. You do not
                  create it by force.&rdquo;
                </Text>
                <Text className="m-0 mt-3 font-12 font-sans text-ink-faint">
                  From a live session on The Creative Act, page 42
                </Text>
              </Section>

              {/* Footer */}
              <Section className="border-border border-t mobile:px-0 px-10 mobile:py-0 py-12">
                <Text className="mt-8 font-12 font-sans text-ink-faint">
                  123 Market Street, Floor 1
                  <br />
                  Tech City, IN, 941026
                </Text>

                <Text className="m-0 mt-4 max-w-[320px] font-12 font-sans text-ink-faint">
                  <Link
                    href={`${baseUrl}/unsubscribe`}
                    className="text-ink-faint underline"
                  >
                    Unsubscribe
                  </Link>{" "}
                  from {brand} product emails.
                </Text>
              </Section>
            </Section>
          </Container>
        </Body>
      </Html>
    </Tailwind>
  );
}

WelcomeEmail.PreviewProps = {
  userFirstName: "John",
  baseUrl: "https://example.com",
} satisfies WelcomeEmailProps;
