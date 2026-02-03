/**
 * Copied from shadc/ui on 12/13/2024
 * @see {@link https://ui.shadcn.com/docs/components/breadcrumb}
 */
import { Slot } from "@radix-ui/react-slot";
import { MoreHorizontal } from "lucide-react";
import {
	type ComponentProps,
	type ComponentPropsWithoutRef,
	type FC,
	type ReactNode,
	type Ref,
} from "react";
import { cn } from "utils/cn";

export const Breadcrumb = ({
	ref,
	...props
}: ComponentPropsWithoutRef<"nav"> & {
	separator?: ReactNode;
	ref?: Ref<HTMLElement>;
}) => <nav ref={ref} aria-label="breadcrumb" {...props} />;
Breadcrumb.displayName = "Breadcrumb";

export const BreadcrumbList = ({
	className,
	ref,
	...props
}: ComponentPropsWithoutRef<"ol"> & {
	ref?: Ref<HTMLOListElement>;
}) => (
	<ol
		ref={ref}
		className={cn(
			"flex flex-wrap items-center text-sm pl-6 my-4 gap-1.5 break-words font-medium list-none sm:gap-2.5",
			className,
		)}
		{...props}
	/>
);

export const BreadcrumbItem = ({
	className,
	ref,
	...props
}: ComponentPropsWithoutRef<"li"> & {
	ref?: Ref<HTMLLIElement>;
}) => (
	<li
		ref={ref}
		className={cn(
			"inline-flex items-center gap-1.5 text-content-secondary",
			className,
		)}
		{...props}
	/>
);

export const BreadcrumbLink = ({
	asChild,
	className,
	ref,
	...props
}: ComponentPropsWithoutRef<"a"> & {
	asChild?: boolean;
	ref?: Ref<HTMLAnchorElement>;
}) => {
	const Comp = asChild ? Slot : "a";

	return (
		<Comp
			ref={ref}
			className={cn(
				"transition-colors font-normal text-content-link hover:text-content-primary",
				className,
			)}
			{...props}
		/>
	);
};

export const BreadcrumbPage = ({
	className,
	ref,
	...props
}: ComponentPropsWithoutRef<"span"> & {
	ref?: Ref<HTMLSpanElement>;
}) => (
	<span
		ref={ref}
		aria-current="page"
		className={cn("flex items-center gap-2 text-content-secondary", className)}
		{...props}
	/>
);

export const BreadcrumbSeparator: FC<ComponentProps<"li">> = ({
	children,
	className,
	...props
}) => (
	<li
		role="presentation"
		aria-hidden="true"
		className={cn(
			"text-content-disabled [&>svg]:w-3.5 [&>svg]:h-3.5",
			className,
		)}
		{...props}
	>
		/
	</li>
);

export const BreadcrumbEllipsis: FC<ComponentProps<"span">> = ({
	className,
	...props
}) => (
	<span
		role="presentation"
		aria-hidden="true"
		className={cn("flex h-9 w-9 items-center justify-center", className)}
		{...props}
	>
		<MoreHorizontal className="h-4 w-4" />
		<span className="sr-only">More</span>
	</span>
);
