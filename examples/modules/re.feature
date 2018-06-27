Feature: Regular Expression
  Scenario: Import regular expression module
    Given a file named "main.cloe" with:
    """
    (import "re")
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly ""

  Scenario: Match pattern with strings
    Given a file named "main.cloe" with:
    """
    (import "re")

    (seq!
      (print (re.match "ab*c" "ac"))
      (print (re.match "ab+c" "ac"))
      (print (re.match "ab*c" "abbbbbc")))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    true
    false
    true
    """

  Scenario: Find pattern in strings
    Given a file named "main.cloe" with:
    """
    (import "re")

    (seq!
      (print (re.find "abc" "ac"))
      (print (re.find "ab*c" "ac"))
      (print (re.find "a(bc+)" "abcccc"))
      (print (re.find "a(b*c)" "abbbbc")))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    nil
    ["ac"]
    ["abcccc" "bcccc"]
    ["abbbbc" "bbbbc"]
    """

  Scenario: Replace pattern in strings
    Given a file named "main.cloe" with:
    """
    (import "re")

    (seq!
      (print (re.replace "abc"    "foo"    "ac"))
      (print (re.replace "ab*c"   "x${0}z" "ac"))
      (print (re.replace "a(bc+)" "x${1}z" "abcccc"))
      (print (re.replace "a(b*c)" "x${1}z" "abbbbc")))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    ac
    xacz
    xbccccz
    xbbbbcz
    """
