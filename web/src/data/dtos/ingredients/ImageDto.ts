export class ImageDto {
    urls: ImageURLs = new ImageURLs();
    author: ImageAuthor = new ImageAuthor();

    constructor(obj: Partial<ImageDto>) {
        Object.assign(this, obj);
    }
}

class ImageURLs {
    small: string = '';
    regular: string = '';
    raw: string = '';
}

class ImageAuthor {
    name: string = '';
    username: string = '';
    profileURL: string = '';
}
