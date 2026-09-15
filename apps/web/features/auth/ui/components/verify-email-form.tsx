"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";

import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRightIcon, LoaderCircleIcon, RefreshCwIcon } from "lucide-react";
import { Controller, useForm } from "react-hook-form";
import type { z } from "zod";

import { Button } from "@loreline/ui/components/button";
import { Field, FieldError, FieldLabel } from "@loreline/ui/components/field";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from "@loreline/ui/components/input-otp";
import { toast } from "@loreline/ui/components/toast";

import { getApiErrorCode } from "@/api/utils";

import { useResendEmailVerification, useVerifyEmail } from "../../api";
import { verifyEmailSchema } from "../../schemas";

const RESEND_COOLDOWN_SECONDS = 30;

type VerifyEmailFormProps = {
  email: string;
};

export const VerifyEmailForm = ({ email }: VerifyEmailFormProps) => {
  const router = useRouter();
  const { mutate: verifyEmail, isPending: isVerifying } = useVerifyEmail();
  const { mutate: resendCode, isPending: isResending } =
    useResendEmailVerification();

  const [cooldown, setCooldown] = useState(0);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const form = useForm<z.infer<typeof verifyEmailSchema>>({
    resolver: zodResolver(verifyEmailSchema),
    defaultValues: {
      email,
      code: "",
    },
  });

  useEffect(() => {
    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, []);

  function startCooldown() {
    setCooldown(RESEND_COOLDOWN_SECONDS);
    intervalRef.current = setInterval(() => {
      setCooldown((prev) => {
        if (prev <= 1) {
          if (intervalRef.current) clearInterval(intervalRef.current);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
  }

  function onSubmit(data: z.infer<typeof verifyEmailSchema>) {
    verifyEmail(
      { body: data },
      {
        onSuccess: () => {
          toast.add({
            title: "Email verified successfully",
            type: "success",
          });
          router.push("/sign-in");
        },
        onError: (error) => {
          const code = getApiErrorCode(error);
          if (code === "INVALID_CODE" || code === "CODE_EXPIRED") {
            form.setError("code", {
              type: "server",
              message:
                code === "CODE_EXPIRED"
                  ? "This code has expired. Request a new one"
                  : "That code doesn't look right",
            });
          }
        },
      },
    );
  }

  function handleResend() {
    if (cooldown > 0 || isResending) return;

    resendCode(
      { body: { email } },
      {
        onSuccess: () => {
          toast.add({
            title: "Verification code sent",
            type: "success",
          });
          startCooldown();
        },
      },
    );
  }

  return (
    <div className="m-auto w-full max-w-md py-14">
      <p className="font-semibold text-brand-ink text-sm">Almost there</p>
      <h2 className="mt-3 font-semibold text-5xl leading-[1.08] tracking-[-0.047em]">
        Verify your email.
      </h2>
      <p className="mt-3 text-muted-foreground text-sm leading-relaxed">
        We sent a 6-digit code to{" "}
        <span className="font-medium text-foreground">{email}</span>. Enter it
        below to confirm it&apos;s you.
      </p>

      <div className="mt-9 space-y-5">
        <form
          id="form-verify-email"
          onSubmit={form.handleSubmit(onSubmit)}
          className="space-y-5"
        >
          <Controller
            name="code"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <div className="flex items-center justify-between">
                  <FieldLabel htmlFor="form-verifiy-email-verification-code">
                    Verification code
                  </FieldLabel>
                  <Button
                    type="button"
                    variant="outline"
                    size="xs"
                    onClick={handleResend}
                    disabled={cooldown > 0 || isResending}
                  >
                    <RefreshCwIcon
                      className={isResending ? "animate-spin" : ""}
                    />
                    {cooldown > 0 ? `Resend in ${cooldown}s` : "Resend Code"}
                  </Button>
                </div>
                <InputOTP
                  {...field}
                  maxLength={6}
                  autoFocus
                  id="form-verifiy-email-verification-code"
                >
                  <InputOTPGroup className="*:data-[slot=input-otp-slot]:h-12 *:data-[slot=input-otp-slot]:w-17 *:data-[slot=input-otp-slot]:text-xl">
                    <InputOTPSlot index={0} />
                    <InputOTPSlot index={1} />
                    <InputOTPSlot index={2} />
                  </InputOTPGroup>
                  <InputOTPSeparator className="mx-2" />
                  <InputOTPGroup className="*:data-[slot=input-otp-slot]:h-12 *:data-[slot=input-otp-slot]:w-17 *:data-[slot=input-otp-slot]:text-xl">
                    <InputOTPSlot index={3} />
                    <InputOTPSlot index={4} />
                    <InputOTPSlot index={5} />
                  </InputOTPGroup>
                </InputOTP>
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />
          <Button
            type="submit"
            size="xl"
            className="w-full"
            disabled={isVerifying}
          >
            {isVerifying ? <LoaderCircleIcon className="animate-spin" /> : null}
            Verify
            <ArrowRightIcon data-icon="inline-end" />
          </Button>
        </form>
      </div>

      <p className="mt-10 text-center text-muted-foreground text-xs">
        Trouble verifying?{" "}
        <Link href="/support" className="underline underline-offset-1">
          Contact support
        </Link>
        .
      </p>
    </div>
  );
};
