import { useTheme } from "@emotion/react";
import { type ImgHTMLAttributes, type Ref } from "react";
import { getExternalImageStylesFromUrl } from "theme/externalImages";

export const ExternalImage = ({
	ref,
	...props
}: ImgHTMLAttributes<HTMLImageElement> & {
	ref?: Ref<HTMLImageElement>;
}) => {
	const theme = useTheme();

	return (
		// biome-ignore lint/a11y/useAltText: alt should be passed in as a prop
		<img
			ref={ref}
			css={getExternalImageStylesFromUrl(theme.externalImages, props.src)}
			{...props}
		/>
	);
};
