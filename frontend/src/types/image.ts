export interface ImageEditorParams {
  RedBrightness: number;
  GreenBrightness: number;
  BlueBrightness: number;
  Contrast: number;
  Negative: boolean;
  Order: string;
  VerticalMirror: boolean;
  HorizontalMirror: boolean;
  Magic: number;
}

export interface ChannelImage {
  type: 'red' | 'green' | 'blue' | 'gray';
  url: string;
}