from tkinter import *
from .tester import *


def handler():
    try:
        limit = int(entry_limit.get())
        prefix = int(entry_prefix.get())
        _, set_ip = get_iplist(prefix, limit)
        output.delete("0.0", "end")
        [output.insert("0.0", item + " status code from server = " + str(make_req(item)) + '\n') for item in set_ip]
    except ValueError:
        output.delete("0.0", "end")
        output.insert("0.0", "Input value")

def handler_random_ip():
    ip = random_ip()
    output.delete("0.0", "end")
    output.insert("0.0", str(ip) + '\n')
    status = make_req(ip)
    output.insert("0.0", "status code from server = " + str(status) + '\n')


root = Tk()

label1 = Label(root, text="Client for ip request")
label1.grid()

frame = Frame(root)
frame.grid()

label2 = Label(frame, text="amount of IP's:")
label2.grid(row=1, column=1)
# limit
entry_limit = Entry(frame, width=10, borderwidth=2)
entry_limit.grid(row=1, column=2)
# prefix
label3 = Label(frame, text='prefix:')
label3.grid(row=1, column=3)
entry_prefix = Entry(frame, width=10, borderwidth=2)
entry_prefix.grid(row=1, column=4)
# Generate
button1 = Button(frame, text="Generate", command=handler)
button1.grid(row=1, column=7, padx=(10, 0))

# limit
label5 = Label(frame, text="Generate random IP and make request:")
label5.grid(row=2, column=1)
button2 = Button(frame, text="Generate and make", command=handler_random_ip)
button2.grid(row=2, column=2, padx=(15, 0))

output = Text(frame, bg="lightblue", font="Arial 9", width=100, height=20)
output.grid(row=3, columnspan=14)

root.mainloop()
