const baseUrl = '/api/b2b';
const acceptInvitation = `${baseUrl}/accept-invitation`;


class Endpoint {
  // post
  readonly login = `${baseUrl}/login`;
  // post
  readonly signUp = `${baseUrl}/signup`;

  // get
  verifyEmail(token: string): string {
    return `${baseUrl}/verify/${token}`;
  }

  // post
  readonly passwordReset = `${baseUrl}/password-reset`;
  // post
  readonly pwResetEmail = `${this.passwordReset}/letter`;
  // get
  verifyPwToken(token: string): string {
    return `${this.passwordReset}/token`;
  }

  // Get
  readonly account = `${baseUrl}/account`;
  // Get
  readonly refreshJWT = `${this.account}/jwt`;
  // Get
  readonly profile = `${this.account}/profile`;
  // Post
  readonly requestVerification = `${this.account}/request-verification`;
  //Patch
  readonly displayName = `${this.account}/display-name`;
  // Patch
  readonly changePassword = `${this.account}/password`;

  // get, post, patch
  readonly team = `${baseUrl}/team`;
  // get
  readonly members = `${this.team}/members`

  // delete
  deleteMember(id: string): string {
    return `${this.team}/${id}`;
  }

  // get
  readonly products = `${baseUrl}/products`;

  // get
  readonly licences = `${baseUrl}/licences`
  // patch, delete
  licenceOf(id: string): string {
    return `${this.licences}/${id}`;
  }

  // get, post
  readonly invitations = `${baseUrl}/invitations`;
  // delete
  invitationOf(id: string): string {
    return `${this.invitations}/${id}`
  }

  // get
  verifyInvitation(token: string): string {
    return `${acceptInvitation}/verify/${token}`;
  }

  // get
  readonly verifyLicence = `${acceptInvitation}/licence`;
  // post
  readonly readerSignUp = `${acceptInvitation}/signup`;
  // post
  readonly granceLicence = `${acceptInvitation}/grant`;
}

export const apiUrl = new Endpoint();
