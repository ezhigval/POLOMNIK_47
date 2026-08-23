"use client";

import Link from "next/link";
import type { LegalDocumentType } from "@/lib/operator-config";
import { legalDocumentPaths } from "@/lib/operator-config";

type ConsentCheckboxProps = {
  name: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
  required?: boolean;
  disabled?: boolean;
  documentType: LegalDocumentType;
  documentTitle: string;
  labelPrefix: string;
};

export function ConsentCheckbox({
  name,
  checked,
  onChange,
  required = false,
  disabled = false,
  documentType,
  documentTitle,
  labelPrefix,
}: ConsentCheckboxProps) {
  const href = legalDocumentPaths[documentType];

  return (
    <label className="flex items-start gap-3 text-sm leading-6 text-stone-700">
      <input
        type="checkbox"
        name={name}
        checked={checked}
        onChange={(event) => onChange(event.target.checked)}
        required={required}
        disabled={disabled}
        className="mt-1 size-4 shrink-0 rounded border-stone-300 text-brand-800 focus:ring-brand-700"
      />
      <span>
        {labelPrefix}{" "}
        <Link href={href} target="_blank" className="font-medium text-brand-800 underline underline-offset-2 hover:text-brand-900">
          {documentTitle}
        </Link>
        .
      </span>
    </label>
  );
}

type PersonalDataConsentProps = {
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
};

export function PersonalDataConsentCheckbox({ checked, onChange, disabled }: PersonalDataConsentProps) {
  return (
    <ConsentCheckbox
      name="consent_personal_data"
      checked={checked}
      onChange={onChange}
      required
      disabled={disabled}
      documentType="personal_data"
      documentTitle="Согласие на обработку персональных данных"
      labelPrefix="Я даю согласие на обработку моих персональных данных в соответствии с"
    />
  );
}

type MarketingConsentProps = {
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
};

export function MarketingConsentCheckbox({ checked, onChange, disabled }: MarketingConsentProps) {
  return (
    <ConsentCheckbox
      name="consent_marketing"
      checked={checked}
      onChange={onChange}
      disabled={disabled}
      documentType="marketing"
      documentTitle="Согласие на получение рекламных и информационно-маркетинговых сообщений"
      labelPrefix="Я хочу получать информационные и рекламные сообщения о туристических, паломнических и иных услугах на условиях"
    />
  );
}

type TermsConsentProps = {
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
};

export function TermsConsentCheckbox({ checked, onChange, disabled }: TermsConsentProps) {
  return (
    <ConsentCheckbox
      name="consent_terms"
      checked={checked}
      onChange={onChange}
      required
      disabled={disabled}
      documentType="terms"
      documentTitle="Пользовательского соглашения"
      labelPrefix="Я принимаю условия"
    />
  );
}

type DistributionConsentProps = {
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
  required?: boolean;
};

export function DistributionConsentCheckbox({
  checked,
  onChange,
  disabled,
  required = false,
}: DistributionConsentProps) {
  return (
    <ConsentCheckbox
      name="consent_distribution"
      checked={checked}
      onChange={onChange}
      required={required}
      disabled={disabled}
      documentType="distribution"
      documentTitle="Согласия на обработку персональных данных, разрешённых для распространения"
      labelPrefix="Я даю согласие на обработку моих персональных данных, разрешённых для распространения, на условиях"
    />
  );
}

type PhotoDistributionConsentProps = {
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
};

export function PhotoDistributionConsentCheckbox({
  checked,
  onChange,
  disabled,
}: PhotoDistributionConsentProps) {
  return (
    <label className="flex items-start gap-3 text-sm leading-6 text-stone-700">
      <input
        type="checkbox"
        name="consent_photo_distribution"
        checked={checked}
        onChange={(event) => onChange(event.target.checked)}
        disabled={disabled}
        className="mt-1 size-4 shrink-0 rounded border-stone-300 text-brand-800 focus:ring-brand-700"
      />
      <span>
        Я разрешаю публикацию предоставленной фотографии и связанных с ней персональных данных на сайте и в
        официальных социальных сетях оператора на условиях{" "}
        <Link
          href="/legal/distribution-consent"
          target="_blank"
          className="font-medium text-brand-800 underline underline-offset-2 hover:text-brand-900"
        >
          Согласия на обработку персональных данных, разрешённых для распространения
        </Link>
        .
      </span>
    </label>
  );
}
