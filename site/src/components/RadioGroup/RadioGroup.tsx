/**
 * Copied from shadc/ui on 04/04/2025
 * @see {@link https://ui.shadcn.com/docs/components/radio-group}
 */
import * as RadioGroupPrimitive from "@radix-ui/react-radio-group";

import type * as React from "react";
import { cn } from "utils/cn";

export const RadioGroup = ({
	className,
	ref,
	...props
}: React.ComponentPropsWithRef<typeof RadioGroupPrimitive.Root>) => {
	return (
		<RadioGroupPrimitive.Root
			className={cn("grid gap-2", className)}
			{...props}
			ref={ref}
		/>
	);
};
RadioGroup.displayName = RadioGroupPrimitive.Root.displayName;

export const RadioGroupItem = ({
	className,
	ref,
	...props
}: React.ComponentPropsWithRef<typeof RadioGroupPrimitive.Item>) => {
	return (
		<RadioGroupPrimitive.Item
			ref={ref}
			className={cn(
				`relative aspect-square h-4 w-4 rounded-full border border-solid border-border text-content-primary bg-surface-primary
			focus:outline-none focus-visible:ring-2 focus-visible:ring-content-link
			focus-visible:ring-offset-4 focus-visible:ring-offset-surface-primary
			disabled:cursor-not-allowed disabled:opacity-50`,
				className,
			)}
			{...props}
		>
			<RadioGroupPrimitive.Indicator className="flex items-center justify-center">
				<div className="size-1.5 rounded-full bg-current" />
			</RadioGroupPrimitive.Indicator>
		</RadioGroupPrimitive.Item>
	);
};
