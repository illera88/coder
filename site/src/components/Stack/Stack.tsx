import type { CSSObject } from "@emotion/react";
import type { Ref } from "react";

/**
 * @deprecated Stack component is deprecated. Use Tailwind flex utilities instead.
 */
type StackProps = {
	className?: string;
	direction?: "column" | "row";
	spacing?: number;
	alignItems?: CSSObject["alignItems"];
	justifyContent?: CSSObject["justifyContent"];
	wrap?: CSSObject["flexWrap"];
	children?: React.ReactNode;
	ref?: Ref<HTMLDivElement>;
} & React.HTMLProps<HTMLDivElement>;

/**
 * @deprecated Stack component is deprecated. Use Tailwind flex utilities instead.
 */
export const Stack = (props: StackProps) => {
	const {
		children,
		direction = "column",
		spacing = 2,
		alignItems,
		justifyContent,
		wrap,
		ref,
		...divProps
	} = props;

	return (
		<div
			{...divProps}
			ref={ref}
			css={{
				display: "flex",
				flexDirection: direction,
				gap: spacing * 8,
				alignItems: alignItems,
				justifyContent: justifyContent,
				flexWrap: wrap,
				maxWidth: "100%",
			}}
		>
			{children}
		</div>
	);
};
