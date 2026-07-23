-- "Before User Created" Auth Hook.
--
-- This function runs INSIDE Supabase for every signup attempt, regardless of
-- how the request arrives (our Go backend, the browser SDK, or a raw curl
-- call using the public anon key). That makes it the only backstop that
-- can't be skipped by bypassing our backend entirely.
--
-- Scope is deliberately narrow: only checks that can't be enforced anywhere
-- else, plus terms-of-service acceptance (a compliance-critical check worth
-- enforcing at every layer, not just the frontend/backend). It intentionally
-- does NOT re-check email format or password complexity:
--   - Email format is already validated by our backend (Gin's built-in
--     validator) for our own app's flow.
--   - Password complexity can't be enforced here at all: Supabase does not
--     document the "before user created" hook payload as including the
--     plaintext password. The only reliable, uncircumventable place to
--     enforce password complexity is Supabase's own project-level password
--     policy (Dashboard -> Authentication -> Policies -> Password
--     Requirements) -- set minimum length 8 and require
--     lowercase/uppercase/digit/symbol there. That is a manual step; it
--     can't be expressed as a file in this repo.
--
-- IMPORTANT: verify the payload column/field names below against Supabase's
-- current "Before User Created" hook docs before applying this in the
-- dashboard SQL editor -- the schema wasn't fully confirmed at the time this
-- was written.
--
-- Setup (manual, requires dashboard/CLI access to the Supabase project):
--   1. Run this file in the Supabase SQL editor (or as a migration) to
--      create the function.
--   2. Dashboard -> Authentication -> Hooks -> "Before User Created" ->
--      point it at public.before_user_created.
--   3. If using the Supabase CLI, add to supabase/config.toml instead:
--        [auth.hook.before_user_created]
--        enabled = true
--        uri = "pg-functions://postgres/public/before_user_created"

create or replace function public.before_user_created(event jsonb)
returns jsonb
language plpgsql
as $$
declare
  full_name text;
begin
  full_name := event -> 'user_metadata' ->> 'full_name';

  if full_name is null or length(trim(full_name)) = 0 then
    return jsonb_build_object(
      'error', jsonb_build_object(
        'http_code', 400,
        'message', 'full_name is required'
      )
    );
  end if;

  if (event -> 'user_metadata' ->> 'terms_accepted') is distinct from 'true' then
    return jsonb_build_object(
      'error', jsonb_build_object(
        'http_code', 400,
        'message', 'you must agree to the terms of service and privacy policy'
      )
    );
  end if;

  -- Optional: block a small set of disposable-email domains. Extend this
  -- list as needed; keep it short and low-maintenance rather than trying to
  -- track every disposable-email provider here.
  if (event ->> 'email') ~* '@(mailinator\.com|10minutemail\.com|guerrillamail\.com)$' then
    return jsonb_build_object(
      'error', jsonb_build_object(
        'http_code', 400,
        'message', 'disposable email addresses are not allowed'
      )
    );
  end if;

  return jsonb_build_object();
end;
$$;
