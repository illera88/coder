/**
 * Copied from shadc/ui on 11/13/2024
 * @see {@link https://ui.shadcn.com/docs/components/table}
 */
import { cva, type VariantProps } from "class-variance-authority";
import type * as React from "react";
import { cn } from "utils/cn";

export const Table = ({
	className,
	ref,
	...props
}: React.HTMLAttributes<HTMLTableElement> & {
	ref?: React.Ref<HTMLTableElement>;
}) => (
	<div className="relative w-full overflow-auto">
		<table
			ref={ref}
			className={cn(
				"w-full caption-bottom text-xs font-medium text-content-secondary border-separate border-spacing-0",
				className,
			)}
			{...props}
		/>
	</div>
);
Table.displayName = "Table";

export const TableHeader = ({
	className,
	ref,
	...props
}: React.HTMLAttributes<HTMLTableSectionElement> & {
	ref?: React.Ref<HTMLTableSectionElement>;
}) => (
	<thead ref={ref} className={cn("[&_td]:border-none", className)} {...props} />
);
TableHeader.displayName = "TableHeader";

export const TableBody = ({
	className,
	ref,
	...props
}: React.HTMLAttributes<HTMLTableSectionElement> & {
	ref?: React.Ref<HTMLTableSectionElement>;
}) => (
	<tbody
		ref={ref}
		className={cn(
			"[&>tr:first-of-type>td]:border-t [&>tr>td:first-of-type]:border-l",
			"[&>tr:last-child>td]:border-b [&>tr>td:last-child]:border-r",
			"[&>tr:first-of-type>td:first-of-type]:rounded-tl-md [&>tr:first-of-type>td:last-child]:rounded-tr-md",
			"[&>tr:last-child>td:first-of-type]:rounded-bl-md [&>tr:last-child>td:last-child]:rounded-br-md",
			className,
		)}
		{...props}
	/>
);
TableBody.displayName = "TableBody";

export const TableFooter = ({
	className,
	ref,
	...props
}: React.HTMLAttributes<HTMLTableSectionElement> & {
	ref?: React.Ref<HTMLTableSectionElement>;
}) => (
	<tfoot
		ref={ref}
		className={cn(
			"border-t bg-muted/50 font-medium [&>tr]:last:border-b-0",
			className,
		)}
		{...props}
	/>
);
TableFooter.displayName = "TableFooter";

const tableRowVariants = cva("", {
	variants: {
		hover: {
			true: "hover:bg-surface-secondary data-[state=selected]:hover:bg-muted",
			false: "",
		},
	},
	defaultVariants: {
		hover: false,
	},
});

interface TableRowProps
	extends React.HTMLAttributes<HTMLTableRowElement>,
		VariantProps<typeof tableRowVariants> {}

export const TableRow = ({
	className,
	hover,
	ref,
	...props
}: TableRowProps & {
	ref?: React.Ref<HTMLTableRowElement>;
}) => (
	<tr
		ref={ref}
		className={cn(
			"border-0 border-b border-solid border-border transition-colors",
			"data-[state=selected]:bg-muted",
			tableRowVariants({ hover }),
			className,
		)}
		{...props}
	/>
);
TableRow.displayName = "TableRow";

export const TableHead = ({
	className,
	ref,
	...props
}: React.ThHTMLAttributes<HTMLTableCellElement> & {
	ref?: React.Ref<HTMLTableCellElement>;
}) => (
	<th
		ref={ref}
		className={cn(
			"p-3 text-left align-middle font-semibold",
			"[&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]",
			className,
		)}
		{...props}
	/>
);
TableHead.displayName = "TableHead";

export const TableCell = ({
	className,
	ref,
	...props
}: React.TdHTMLAttributes<HTMLTableCellElement> & {
	ref?: React.Ref<HTMLTableCellElement>;
}) => (
	<td
		ref={ref}
		className={cn(
			"border-0 border-t border-border border-solid",
			"p-3 align-middle [&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]",
			className,
		)}
		{...props}
	/>
);
TableCell.displayName = "TableCell";

const TableCaption = ({
	className,
	ref,
	...props
}: React.HTMLAttributes<HTMLTableCaptionElement> & {
	ref?: React.Ref<HTMLTableCaptionElement>;
}) => (
	<caption
		ref={ref}
		className={cn("mt-4 text-sm text-muted-foreground", className)}
		{...props}
	/>
);
TableCaption.displayName = "TableCaption";
