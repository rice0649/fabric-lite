See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.


Search in video
Intro
0:00
So you want to use Claude Code. You want to get  the most of it, but you don't know exactly how.   This is a crash course, how to master Claude  Code, and we explain it in the most simple  
0:10
way. There are thousands, literally thousands  of other Claude Code tutorials on the internet,  
0:16
but there are none as simple as this. I  brought on Professor Ross Mike. He comes   on and he shares it in the simplest way so that  anyone could create jaw-dropping startups and  
0:26
software using Claude code. We're going to give  you the exact steps for how you can set it up,  
0:32
thinking about the beginner, how to think about  the terminal, how to think about prompting.   But if you stick around to the end of this  episode, there's a tips and tricks section,  
0:41
which I think is super valuable. And  I can't wait to see what you build.  
0:54
we got ross mike on the pod by the end of  this episode what are people going to learn  
1:00
hopefully you're going to not feel overwhelmed  with cloud code i know the terminal is scary   and it's a big boogeyman but i'm going to give  you the blueprint how to use it i'm also going  
1:09
to share consider this the ultimate crash course  on how to use cloud code or any agent effectively  
1:16
Okay, let's get into it. So, I mean, the best  way to start these episodes is with sharing our  
Claude Code Best Practices
1:22
screen. So, when we think of building applications  using AI, using some sort of agent like CloudCode  
1:30
or OpenCode or Codex, whatever it is, there's  a couple of things that you always have to keep  
1:35
in mind. You know, the principles never really  change. One thing that it's important for us to   understand is however good your inputs are will  dictate how good your output is. We're getting to  
1:46
a point where the models are so freakishly good  that if you are producing quote-unquote slop,  
1:54
it's because you've given it slop. There was a  time where the models weren't good enough. There   was a time where we had serious qualms and issues  with the quality of code the models gave us. But  
2:05
now we're starting to get to a point where even  myself, like I'm reviewing a lot more code than I  
2:11
write. And I never thought I'd be able to say that  in the early months of 2026. So very important for  
2:18
us to understand our inputs, how good they are,  how precise they are, how articulate they are,  
2:23
are just as good as our outputs and will dictate  just how good our outputs will be. And the way I  
2:28
want people to think about this is, Greg, is like  imagine you were communicating this to a human,   to a human engineer, right? If you give them  sparse instructions and if anyone is in client  
2:40
work, you realize that most clients, they tell  you one thing, but you have to sort of extract  
2:45
the deeper thoughts of what it is they want. It's  the same way when we work with these agents, when  
2:51
we work with Cloud Code. We need to be really,  really precise with how we build our inputs. Now,  
2:58
what do I mean by inputs? What I mean is our  PRDs or our to-do list or our plans, right? Like  
3:05
there's – people are giving you different names.  It doesn't really matter. It's all the same thing,  
3:10
right? And when we think of a PRD or when we think  of a to-do list or when we think of a plan, I want  
3:15
us to think in such a way as this. Let's say I'm  trying to build this product, right? Let's say  
3:22
– I don't know, Greg, any product ideas that –  Do you have product ideas? Yeah, that's actually  
3:28
the best person to ask, right? Let's say I go on  IdealBrowser.com. I was just going to it. I was  
3:38
just going to it. Yeah, pick the idea of the day  from Idea Browser. It says it's a diagnostic tool  
3:44
for appliance techs losing hundreds of repeat  visits. See, I have no idea what that means,   but let's say I know what that means. Essentially,  when thinking of this idea and looking to build  
3:54
this into a full-fledged product, Generally, the  way you're going to think is, okay, if product X  
4:01
does Y, Z, A, B, and C, how I would reach that is  I'm going to think of features, right? So let's  
4:08
say there's four core features to this application  that Greg just mentioned. And if I have these four  
4:15
features built out, we can safely assume that  we have said product, right? The way we are to  
4:22
design our PRDs to-do lists and plans is such that  we want the agent, the model, to build out all  
4:29
these features, right? Because all these features  put together is our product. You see, a lot of  
4:35
times people will describe a product, not describe  features, and will be frustrated with AI. Like AI  
4:41
is supposed to magically know what you're thinking  about. By the way, Greg, am I making sense so far?  
4:46
100%. I'm with you. Yeah. So we really need to  think in features. But here's the cool part. When  
4:53
developing features, oftentimes the issue with  models is like you'll develop a feature or like  
4:59
let's say the model develops a feature. We don't  know if it works. We don't know if it did it the   right way. That's where with all the cool Ralph  stuff that's happening, we can introduce tests.  
5:08
Right. So let's say the model, the agent builds  feature one. Before moving on to feature two,  
5:15
what I'm going to do is I'm going to get the  model to write a test. If that test passes,   then we'll work on the second feature. If that  test passes, we work on the third feature.  
5:24
So we're finally entering an era where you can  really build something serious with these models.  
Claude Code Plan Mode
5:32
So instead of telling you about just planning,  why don't we do actual planning together? So I'm  
5:38
going to pop up my terminal. So I know everyone's  afraid of the terminal, but in all honesty, if you  
5:43
don't know how to use a terminal, ask AI. It's the  simplest thing. And if not, you can even download  
5:49
the Cloud Code app and go on Code section, give  it a specific folder you want to work on. and use  
5:55
the app. Like there's literally no excuse to  not use Cloud Code. If you're afraid, boohoo,   just jump into use AI, you have all the tools.  That being said, I'm just going to type in Cloud,  
6:04
and we're going to have Cloud Code open. And  usually how people plan is they'll click shift  
6:10
tab, right? And then you have plan mode on.  And you can say, let's say, I want to build  
6:18
tick tock UGC generating app for my marketing  agency. I see like these UGC apps everywhere.  
6:29
please help me create a plan, write this in the,  in, PRD.md file. So this is how most people have  
6:45
planning set up, right? You'll tell Claude Code or  Cursor or whatever agent to do the plan for you,  
6:52
and you ask it to be in some file. And like it  says, it'd be happy to help you plan this out.  
6:57
And it'll ask you some questions, etc, etc. But  I found that there's a better way to get an even  
7:03
more concise plan. And this way, it actually  gets you to think a lot more about trade-offs,  
7:10
concerns, UI, UX decisions, because most of the  time you're sort of allowing the AI to have free  
7:16
reign over certain decisions, which I think will  lead you with a finished product that you're   not excited about. And that's invoking a special  tool. I was going to show you guys the tweet, but  
7:26
unfortunately Twitter is down right now. But Cloud  Code has a specific tool called Ask User Question  
7:32
Tool. And essentially what this tool does, it  starts to interview you about the specifics of  
7:38
your plan right so I'm gonna drop this prompt  where it says read this plan file interview me  
7:44
in detail using the ask user question tool about  literally anything technical implementation UI  
7:49
UX concerns and trade I spelt implementation wrong  do not judge me And what this is going to do is it  
7:56
going to go past the plan that we have and start  to ask us about my new details. So let's finish  
8:01
off this plan first. I'm just going to accept this  is internal use. TechSack will use React. I just  
8:09
want core features. We'll submit answers. And then  CloudCode, you'll see, might ask us a few more  
8:15
questions, but this will generally be the plan.  Right. So it's not just the plan. It's the right  
8:22
plan, right? Like to what you were saying, like  go back, scroll back up here, the features and,  
8:28
yeah, the features and tests. Like the way I  think about this, and I don't know if you agree,  
8:34
is like if you ask Cloud Code to build you a car,  it doesn't really know what a car is. It doesn't  
8:40
understand like you need a steering wheel and a  radio and you need wheels. So the hard part is  
8:47
trying to figure out, is basically explaining  what those things are in a really succinct   and clear way. And that's what this interview is  basically doing. It's explaining each of them and  
8:58
then we're going to test each of those features.  Exactly. Like think of it this way, like a simple   example. Let's say you ask the AI agent to build  you a specific feature, right? How is it going to  
9:10
present that specific feature? Did you want it  in a dashboard? Did you want it to be a modal?   Did it have to be a separate page? Like  when you don't specify these minute details,  
9:19
it will make the assumption for you. And with  Ralph loops and all these type of things, like   you might have a whole application built out and  it's not exactly to the liking or the expectations  
9:29
you had. Right. So let me continue. I'll just make  some selections here just so we can move on and  
The Ask User Question Tool
9:36
then hit submit. and then I'm going to pause this  planning here and then I'm going to paste this.  
9:43
I'm going to say read this plan file and I'm going  to tag the plan file. It's called prd.md. We have  
9:50
that right here and I'm going to say interview  me the details about this question or I don't  
9:56
even need to tag it because it has it in its  context but I just want to show you how annoyingly  
10:02
annoying this is going to get meaning, it's going  to keep asking me questions about said plan or  
10:08
said app idea. So notice how it says round  one core workflow and technical foundation,  
10:15
right? And some of the questions it might even  ask you are things that you might not know about   because you're not technical. So what do I do  when I don't know something, Greg? I'm going to  
10:22
copy this and I'm going to go to the chatbot of  my choice, whether it's Claw, ChatGVT, whatever,   and I'm going to ask questions. So if you remember  earlier, it asked me generic questions about the  
10:32
app. Now it's saying, what's your ideal workflow  for generating UGC video from start to finish?  
10:38
Like notice how the questions are even more  specific now. So it says linear step-by-step,  
10:43
template-based, batch processing, iterative,  conversational. So let's say I select that  
10:49
and it says, how should the app handle Hagen API  costs and usage? So now it's talking about costs,  
10:55
right? Again, most of the times when you just have  a basic plan. This is not included in the plan,  
11:00
right? Let's say we want to have a hard, hard  budget. What database and hosting approach do   you want to use? Most of you probably watching  this have no idea. So I can copy this over, go  
11:09
to chat GPT and ask what's the best decision. This  is my current situation. And then you keep going,  
11:14
you keep going and you submit answers. So when  you use this ask user question tool, the questions  
11:21
become more granular. So it asked me about core  workflow and technical foundation. Now it's going  
11:26
to ask me about UI, UX, and script generation. If  you notice the first plan that it came up with,  
11:32
the default plan for Cloud Code, it was pretty  basic. Now it's asking me, okay, what AI do I want  
11:39
to use for the script generation? I'll use Cloud.  What UI style aesthetic are you going for? Minimal  
11:44
clean, dashboard heavy, creative tool feel, chat  first, right? So hopefully, Greg, I'm making sense  
11:50
with like how much more questions I'm being asked  when I'm invoking this ask user question tool.  
11:59
Yeah, it makes complete sense. You're also going  to use less tokens in the end, right? Yeah,  
12:06
because the thing is, the better your plan, the  better your input, the better the initial set of  
12:12
documents that you give the model, the better  the outcome. And if the better the outcome,  
12:18
There's no back and forth, right? Most people will  have a Ralph's loop running. It'll be a basic plan   and it'll do what you told it to do, but you  weren't specific. So now you're going back and  
12:27
then maybe you're running another loop or you're  going back and doing all these changes. But if you   get it done right, if you invest the time in the  planning stage, I 100% believe you'll save a lot  
12:38
more money. And this will help you clear up a lot  of ideas. For example, this idea that we just had,  
12:44
this TikTok UGC farm, how do we want it set  up? Do we want it to be flat with search? Do  
12:51
we want it to be client campaign? There's a lot  of like these minute details that you're not  
12:57
thinking about. And because you're not thinking  about it, you're allowing Cloud Code to make   those assumptions for you. Right. Which at the  end, after it's burnt through a ton of tokens,  
13:05
now you're going back to change. Right. We can  save so much headache if we do the proper planning  
13:12
from the beginning. And hopefully people see  value in this ask user question tool. Make sure  
13:20
you specify it in your prompt. And hopefully,  Greg, that made sense. It does. So I would say  
13:27
step number one for this Claude Crash Course is  I would get good at planning. I would get really,  
13:32
really good at planning. I would get good at  generating these. Like, look, it keeps on asking  
13:38
me questions. If you notice the very first plan  that we generated with Claude, it was two sets  
13:43
of questions and it was ready to build. But with  this, it's asking me, do I have basic avatars,   custom avatars, multi-scene videos? How do I  want to handle storage? Do I want to download the  
13:53
videos instantly, cloud storage, external storage?  Like there's so much to software engineering. And  
13:59
I think in our last video, someone shared this  on Twitter. I don't know if it was you or someone   else. like software building personal software is  easy but building software others are going to use  
14:10
is very very difficult and if you don't have the  audacity or the decency to to set up a little time  
14:16
a little extra time to plan then i guarantee you  whatever you generate is going to be ai slop and   you might blame the model but really the problem  is you so invest in your plans spend time using  
14:27
planning. Don't use the generic plan mode that  cursor or cloud code has. I would use cloud code  
14:34
and then I would specify the ask user question  tool. It's going to continue to know you with  
14:40
questions like it keeps asking, right? Because  until it knows exactly what it is you want,  
14:45
it won't start building. So I would say that's  step number one to building with cloud code. Step  
Don’t start with Ralph automation (get reps first)
14:52
number two, and everyone's talking about Ralph and  it's exciting, but I wouldn't use it. I wouldn't  
14:59
use Ralph. And the reason I wouldn't use Ralph if  I was just starting out, Greg, is because how are  
15:06
you going to – like imagine this. Like imagine not  knowing how to drive but then buying a Tesla for  
15:15
like the self-driving stuff. Like, cool in theory,  but maybe it's a great idea to know how to drive,  
15:20
how to steer, how to hit the corners, how to maybe  yell at someone when they cut you off before you  
15:26
get the full automated version. I say this to say  because when you get good at developing plans and  
15:34
then working with the AI to build each feature  and testing each feature you start to develop the  
15:40
sense on product building on like you know even I  heard someone called Vibe QA testing You get this  
15:48
sense by going one on one yourself And this is  why a lot of people who were fighting with Cloud  
15:54
Code all these months are really, really good at  using it now because they spent the time building  
15:59
without using these crazy automation. So if you're  using Cloud Code for the first time or you're  
16:05
just getting into it, good plan, number one. And  number two, get your reps in by not using Ralph.  
16:11
So develop the features one by one. Now that you  have your plan, you can literally tell Cloud Code,  
16:16
hey, OK, let's build the first feature. You know,  go ahead and do it. And then once the feature is  
16:21
done, you can test it or ask it. How can I test  this? How can I run this app? I wouldn't jump into  
16:27
using Ralph right away. Build without Ralph.  But let's say you've built these reps now and  
What are “Ralph loops” and why plans and documentation matter most
16:36
you're comfortable with Cloud Code. Now you hear  about all these things, skills, MCP, prompt.md,  
16:44
agent.md. What else is there? Something.md.  You hear all these conventions, plugins. You  
16:52
have Ralph, all these things. So what do I need to  perfectly build something using CloudCode or any  
16:59
agent? I'll be honest with you. Most of these  things are all the same. PromptMD and AgentMD  
17:06
are just markdown files. Plugins are skills  with a little bit extra. What you need to build  
17:14
successfully using these agents is, first of all,  you need a good plan, right, which are documents,  
17:21
which is the PRD.MD we just generated. And then  you need to document the progress that's being  
17:29
made. For anyone who's familiar with Ralph, you  know what I'm talking about. For those who aren't,  
17:36
what's cool about a Ralph loop is as follows. A  Ralph loop is basically you have a list of things  
17:42
that need to get done. The whatchamacallit, the  PRD.MD or the plan. You give it to the AI model.  
17:51
The model works on the first task. It finishes  and then documents it in another file and then  
17:57
it goes again and it stops until it's completed  the whole list. Now, this isn't anything special,  
18:05
but the reason why it's now super powerful  is because the models are getting so,   so good. But here is the issue. If you have  a terrible plan, if you have a terrible PRD,  
18:15
this doesn't matter. You're just donating money  to Anthropic and I wish you the best of luck if   that's what you want to do. But if you want  to make sure that your tokens are not wasted,  
18:25
you're going to invest in a good PRD.MD file or  a good plan file. Greg, am I making sense so far?  
18:33
A hundred percent. OK. You're driving the point  home. Yes. So I'll talk a little bit about Ralph  
18:40
now. So with Mr. Ralph Wiggum, how do we use this?  Now, there's a lot of different iterations, like  
Ras’s Ralph setup: progress tracking + tests + linting
18:48
people are coming with their own style. I'm going  to share with you my Ralph setup in a second,   Greg. One thing I will say is Cloud Code has a  plug-in, a Ralph Wiggum plug-in. I wouldn't use  
18:59
that. And the reason I wouldn't use that is even  the person who invented the whole Ralph system  
19:04
is against it. It's not the best use of Ralph.  But I just want to share this concept of how  
19:10
Ralph works. It's essentially going to go through  our plan, and it's going to build out each feature  
19:16
step by step. and it's not going to stop until  it's done. This is cool when your plan rocks.  
19:23
If your plan sucks, then it's terrible. It  doesn't matter. Now, in terms of how to set up  
19:31
Ralph Wiggum, I have my own setup and I don't want  anyone to think I'm shilling my own setup for any  
19:37
reason. But the reason why I built my own setup is  there's a couple of things my Ralph loop does. The  
19:44
first thing is it makes sure that there's a plan,  a prd.md file, and there's a progress.txt file.  
19:51
But it also, every feature it builds, it then  writes a test and it then lints. And basically  
19:58
what this does is it makes sure that every feature  that's built actually works, right? Because  
20:03
there's no point on working on feature two if  feature one doesn't work. If feature one doesn't   work, if the test fails, guess what the AI model  is going to do? It's going to go back to working  
20:12
on feature one. And once the test passes, we work  on feature two. And then once feature two test  
20:18
passes, we work on features three. Right. All this  is awesome. But I'm going to go back to the same  
20:24
point. If your plan sucks, then the Ralph loop  won't matter. Now, in order to set up this loop,  
20:32
you can find the get up here how to set it up. You  honestly, I'm not even going to explain it. Greg,  
20:38
people can literally copy the link, pass it  to Claude and then be like I want to run this   Ralph loop and it will tell you exactly what  to do that's how good the models have become  
20:48
but I'll show you an example of this running so  I have a simple prd.md file it's nothing crazy  
20:55
it's just to show you the point but basically  there are a couple tasks here I want to build   a basic server that has some basic endpoints and  I just want to show you how my Ralph loop works.  
21:06
So when I run this Ralph loop, and again, if you  don't know how to run this, you paste the GitHub  
21:12
URL in CloudCode in your agent and ask it and it  will tell you how to do it. I have a few different  
21:18
configurations. I can use open code if I want. I  can use codex if I want, but I'm just going to use   CloudCode and I'm just going to run the script.  And basically what it's going to start doing is  
21:29
it's going to start running through each task,  as you can see, and it's going to update the PRD  
21:36
and it's just going to continue to work. Now, I  can go and leave, right? I can go about my day,  
21:42
hang with Greg and this loop will continue to  work. And I'm going to see that at some point,  
21:50
whether it's five minutes, three minutes, 10  minutes, however long this is, this is going   to finish all the tasks. I'm going to have a  working product built and all this is cool,  
22:00
but it doesn't matter if I'm going to go back to  the original document, if the plan isn't good.  
22:08
Now, skills are great. MCPs are great. All these  different markdown files are great. You would do  
22:13
yourself a serious service if your plan is good.  So the key to successfully building with CloudCode  
22:23
is you have an absolutely great plan. And if  you use the ask user question tool, You will  
22:29
spend so much time on the plan where it starts to  get annoying. It doesn't get fun. But those of us  
22:35
who focus on this will end up having better  outputs. Let's continue. If you notice here,  
22:42
my Ralph loop is continuing to go and it took care  of the first task. I can see some files already  
22:47
generated. If I go to the progress.txt file, you  can see, Greg, it started to make some progress.  
22:54
It's documenting that. And this is just going to  continue to work. This is just going to continue   to run. So people have different iterations. I  know the AmpCode people have their own iteration  
23:03
and different people have their own iteration. It  doesn't really matter. Right. Someone's Ralph is   could be better. Someone's can be worse. Someone's  going to be all of that is cool. But don't get  
23:13
stuck in the weeds. The main sauce is how you can  articulate perfectly in a beautiful presentation.  
23:21
create the perfect input because if you create  the perfect input we have reached a point where   the models will give you perfect output so that  my main tip crash course for people use the ask  
23:34
user question tool build without using Ralph And  if you are going to use Ralph understand if your  
23:40
plan sucks you just donating money to Anthropic  And I think Anthropic has enough money that they   don't need your money being donated to them. Amen.  Amen. Is there anything else people need to know,  
Tips & tricks: don’t obsess over MCP/skills/plugins
23:53
like little tips and tricks? I notice, you  know, you're not using the Mac terminal,  
23:59
you're using Ghosty. ghosty yes yes so honestly  it's all preference right so like the terminal you  
24:05
use and all this stuff is all preference here's  what i would say like let's have a tips and tricks  
24:12
list tips and tricks so first i would say is my  goodness spelling today first i would say is use  
24:22
the ask what was the specific tool i just want to  make sure i don't forget ask user questions tool  
24:28
Slept on. I don't know why no one's not talking  about it. It literally I saw the tweet from the   Anthropic team. Hundred percent. I would use that  when planning. Number two, don't over obsess.  
24:43
Obsess on MCP skills, et cetera, et cetera. I'm  not saying don't get into these. I'm not saying  
24:49
don't read about them. I'm not saying don't use  them. But I can almost guarantee you these things  
24:55
are not the reason why your product isn't working,  right? Most of the time, it's your plan sucks,  
25:02
right? That's number two. Number three, I would  use Ralph after I've built something without. And  
25:12
the reason being is, again, listen, if you are a  baller shot caller and you have all the money to   blow and you don't care and you want to donate  money to Anthropic, go ahead and use Ralph. But  
25:21
If we were to sit here eye to eye and you haven't  built anything, deployed anything, there isn't a  
25:26
URL that I myself or Greg can click on that you've  built. You have no business using Ralph. You  
25:32
literally have no business using Ralph. I would  first get good at prompting and building something  
25:38
using a plan, whether it's whatever AG1, Cloud  Code, Open Code, whatever. Once you have something  
25:44
deployed to Vercel or like there's a URL and we  can use it, then you can use Ralph. So number  
25:51
four, this is a little in the weeds, but context  is more important than ever. And a lot of times,  
25:59
Cloud Code or even Cursor will tell you what  percent of context has been used. I generally  
26:05
wouldn't go over 50%, meaning like the Anthropic  model, Opus 4.5 has a 200,000 token context limit.  
26:13
The moment, in my opinion, you've got over 100,000  tokens, meaning you're using the same session. it  
26:19
starts to sort of deteriorate that's when you have  people greg who say oh like i started off good but  
26:24
it started going bad that's because you've filled  it with so much context and the best way to think  
26:30
about this is like yourself right like let's say  we went to some english class and or some you know  
26:36
whatever class and the professor just kept dumping  information information information at some point  
26:41
we're going to feel overwhelmed and we're going  to actually start forgetting stuff And I'm not   saying that's how the models work, but that's how  the models act. Right. So context is very much  
26:51
important. The moment you see 50 percent or even  40 percent, I would start a new session. And last  
26:57
but not least, have audacity. And what I mean by  that is software development is starting to become  
27:04
easy. The software engineering is very, very hard.  And what do I mean by that? To architect software,  
27:10
to make sure things are usable, to create great  UX, UI. to have great taste, to make something  
27:16
that people actually use requires time. And  in order to spend time, it requires audacity.   I know the models are good and you can clone  a six billion dollar software. But if all of  
27:26
us can do it now, what makes software different?  I think thinking about those things and thinking   about the art of building products and building  something that's tasteful is very, very important.  
27:36
And I think anyone who uses these five tips should  kick cheeks in 2025, 2026. Sorry. I agree on the  
Scroll-stopping software wins
27:46
audacity thing. I think like it's for me, it's  like about creating scroll stopping software.  
27:51
You know what I mean? Like there's so many  people and there's a lot of tutorials about this,   like cloning billion dollar software. You know, I  cloned a four billion dollar software. Look at me.  
28:01
But that's not the type of software that's going  to work in 2026. Right. I saw this. Let me just  
28:09
share it real quick. I saw this guy who created a  running app based on how you're feeling. So it's  
28:15
like, how are you feeling? Stressed, angry, and  it's an AI-assisted running app that interprets  
28:22
your current emotions to generate a personalized  route. And I just thought it was interesting.   You know what I mean? I'd never seen an app like  this. And I think that as you call it Audacity,  
28:33
I think this is an audacious app, right? It's  scroll-stopping. You haven't seen it before. So  
28:39
I think push, you want to push Cloud Code to like  get you to this basically. And this is why I'm  
28:47
like so pro people not using Ralph if they haven't  built anything fully. Because like now we're,  
28:53
people are getting to a point where they want the  model to think for them, right? Where like if you   look at the app, you just shared the animations  and how things were floating and like even the  
29:02
colors used for the different emotions. Like that  requires thought, right? And that's what stops   people now. Like if building the AI chat interface  is easy, what's going to make your app different?  
29:12
I think a little bit of audacity, a little bit  of thought and care and a little bit of taste   goes a long way nowadays. And more than the models  getting better because it's going to get easier.  
29:22
It's going to get better. It's going to get  faster. But unfortunately, if you don't change,   then it doesn't all matter. Yeah. And don't  be afraid to use pen and paper like this. This  
29:32
person literally just like started sketching out  the features. Yeah. Like how did this thing work?  
29:38
Yeah. I love it. I love it. Right. And this  is why the app, if I don't know the metrics,  
29:43
but I'm willing to bet it's doing really, really  well because all this stuff matters. Like we could   clone something like this feature wise, but I'm  willing to bet like the feel, the animations,  
29:53
the colors, we would not be able to get it  exactly like this. 100% all right man thanks  
30:00
for coming on you got me fired up I actually I  didn't know about that interview tool so thanks  
30:06
for sharing that with me yeah just a heads up  it will ask a lot of questions I shared it with  
30:12
a couple friends and a couple people got annoyed  but it's worth it right especially if you wanted  
30:17
to build something end-to-end or you're building  a very like like very minute detailed feature then  
30:24
it's really really worth it i wouldn't use the  general plan personally so just a heads up but  
30:29
it's really really worth it and i would love to  hear people's feedback in the comments sounds good  
30:34
we'll be in the comments you got to come back on  in a few months or whenever people want you it's  
30:41
always an absolute privilege to have you here  i'll include links where you can follow and you  
30:46
should follow monsieur rass mike his youtube  channel is x i'll include the link to ralphie  
30:55
even though if you're a beginner don't even click  that link i i wouldn't like i know there's maybe  
31:01
some degenerates who do but i highly suggest you  don't because if you haven't even built without  
31:06
it then no point have some willpower folks come  on you know don't click the link but i'm putting  
31:12
it in there because i want to see who's tempted  and yeah thanks again for coming on i'll see you  
31:18
i'm coming to toronto in april so let's hang out  we'll see we'll see each other then and again as   always it's a pleasure thank you so much you know  for bringing me on of course later have a good one
