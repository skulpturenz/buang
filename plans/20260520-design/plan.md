# Description

Plan and design a user interface, logo, logo and color scheme for Buang. Buang is a REST API which deploys docker compose services for short lived feature branches while they are open. This facilitates collaboration by reducing friction when other members in a team need to verify changes which are being introduced. Buang means throw away in Malay and it was named so because Docker containers are disposable.

# Scope

There are 4 pages which need to be designed:
1. Login
The login page allows users to login to the system. Users can login via email and password or they can login with SSO. Users should also be able to navigate to the registration page from the login page

2. Register
For users who are not using SSO, they can also register with their email, password, buang api base url and buang api key. These are all required fields.

3. Setup
This page allows users to setup a repository and Buang project. Users must specify:
   - Display a form to capture:
      - The repository URL
         - Required
      - Whether the repository requires authentication
         - Required
      - If authentication is required, the username to authenticate with
         - Conditionally required
      - If authentication is required, the password to authenticate with
        Display helper text which reminds the user to use a Personal Access Token (PAT) instead of their real password
         - Conditionally required
      - Environment variables. This will be a list of key value pairs which are ideally represented by a table which allows rows to be added and deleted
         - Conditional
      - Service entrypoint
         - Required
         - The main service which will accept the http requests. This is usually a load balancer like Traefik or NGINX
      - Wait for workflow run
         - Conditional
         - Some projects will require a build and push workflow to complete before it is ready to deploy. In these cases the deployment
           shouldn't happen immediately, it should only occur after the workflow run is successful

# Related
- Figma files for shadcn components are attached
- The current FE implementation is attached
- We should have both light and dark themes. An example Tailwind theme from another project is:
```css
@import "tailwindcss";
@config "../tailwind.config.js";

@layer base {
	:root {
		--background: 0 0% 100%;
		--foreground: 0 0% 0%;

		--card: 0 0% 98%;
		--card-foreground: 0 0% 0%;

		--popover: 0 0% 99%;
		--popover-foreground: 0 0% 0%;

		--primary: 67 90% 56%;
		--primary-foreground: 0 0% 0%;

		--secondary: 240 1% 90%;
		--secondary-foreground: 0 0% 0%;

		--muted: 0 0% 98%;
		--muted-foreground: 0 0% 45%;

		--accent: 70 86% 95%;
		--accent-foreground: 0 0% 0%;

		--destructive: 0 84% 60%;
		--destructive-foreground: 0 0% 100%;

		--border: 0 0% 95%;
		--input: 0 0% 85%;
		--ring: 0 0% 80%;

		--radius: 0.5rem;
	}

	.dark {
		--background: 0 0% 0%;
		--foreground: 0 0% 100%;

		--card: 0 0% 3%;
		--card-foreground: 0 0% 100%;

		--popover: 0 0% 3%;
		--popover-foreground: 0 0% 100%;

		--primary: 67 90% 56%;
		--primary-foreground: 0 0% 0%;

		--secondary: 0 0% 10%;
		--secondary-foreground: 0 0% 100%;

		--muted: 0 0% 9%;
		--muted-foreground: 0 0% 45%;

		--accent: 0 0% 10%;
		--accent-foreground: 0 0% 100%;

		--destructive: 0 72% 51%;
		--destructive-foreground: 0 0% 100%;

		--border: 0 0% 15%;
		--input: 0 0% 20%;
		--ring: 67 90% 56%;
	}
}

@layer base {
	* {
		@apply border-border;
	}

	body {
		@apply bg-background text-foreground;
	}
}
```

# Questions answered

- deliverable_layout: Design canvas — all 4 pages side-by-side in light + dark for easy comparison
- logo_direction: Mark: branch/Git-flow inspired (forking, ephemeral leaf)
- color_scheme: Decide for me
- aesthetic: Terminal-coded — mono-heavy, grid lines, brutalist accents
- auth_layout: Centered single column (current implementation)
- setup_layout: Two-column — form left, live summary / preview right
- theme_mode: Show both side-by-side for every page
- extras: Empty / success state for 'repo configured', Favicon / app icon variants, A small color/type token sheet, Error state mocks (failed login, API down)
- tone_copy: Tighten current copy where helpful, otherwise same

# Follow up

1. Propose a different color scheme
